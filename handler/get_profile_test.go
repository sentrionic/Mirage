package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/mirage/mocks"
	"github.com/sentrionic/mirage/model/apperrors"
	"github.com/sentrionic/mirage/model/fixture"
	"github.com/sentrionic/mirage/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetProfile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	uid, _ := service.GenerateId()

	t.Run("Success", func(t *testing.T) {
		mockUserResp := fixture.GetMockUser()

		mockUserService := new(mocks.UserService)
		mockUserService.On("FindByUsername", mockUserResp.Username).Return(mockUserResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := gin.Default()
		store := cookie.NewStore([]byte("secret"))
		router.Use(sessions.Sessions("mqk", store))

		router.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("userId", uid)
		})

		router.Use(func(c *gin.Context) {
			c.Set("userId", uid)
		})

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/v1/profiles/"+mockUserResp.Username, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockUserResp.NewProfileResponse(uid))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})

	t.Run("NoContextUser", func(t *testing.T) {
		mockUserResp := fixture.GetMockUser()

		mockUserService := new(mocks.UserService)
		mockUserService.On("FindByUsername", mockUserResp.Username).Return(mockUserResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := gin.Default()
		store := cookie.NewStore([]byte("secret"))
		router.Use(sessions.Sessions("mqk", store))

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/v1/profiles/"+mockUserResp.Username, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockUserResp.NewProfileResponse(""))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		username, _ := service.GenerateId()

		mockUserService := new(mocks.UserService)
		mockUserService.On("FindByUsername", username).Return(nil, fmt.Errorf("some error down call chain"))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/v1/profiles/"+username, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respErr := apperrors.NewNotFound("profile", username)

		respBody, err := json.Marshal(gin.H{
			"error": respErr,
		})
		assert.NoError(t, err)

		assert.Equal(t, respErr.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t) // assert that UserService.Get was called
	})
}
