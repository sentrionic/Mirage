package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/mirage/mocks"
	"github.com/sentrionic/mirage/model"
	"github.com/sentrionic/mirage/model/apperrors"
	"github.com/sentrionic/mirage/model/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := fixture.GetMockUser()
	reqUser := &model.User{
		Email:       user.Email,
		Password:    user.Password,
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mockUserService *mocks.UserService)
		checkResponse func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":    user.Username,
				"password":    user.Password,
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", reqUser).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
				respBody, err := json.Marshal(user.NewAccountResponse())
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
				mockUserService.AssertCalled(t, "Register", reqUser)
				assert.Contains(t, recorder.Header(), "Set-Cookie")
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":    user.Username,
				"password":    user.Password,
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", reqUser).Return(apperrors.NewConflict("User Already Exists", reqUser.Email))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Username too short",
			body: gin.H{
				"username":    "te",
				"password":    user.Password,
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "Username too long",
			body: gin.H{
				"username":    fixture.RandStringRunes(20),
				"password":    user.Password,
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"username":    user.Email,
				"password":    user.Password,
				"displayName": user.DisplayName,
				"email":       "test.email",
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "Password Too Short",
			body: gin.H{
				"username":    user.Email,
				"password":    "pass",
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "DisplayName too short",
			body: gin.H{
				"username":    user.Username,
				"password":    user.Password,
				"displayName": "tes",
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "DisplayName too long",
			body: gin.H{
				"username":    user.Username,
				"password":    user.Password,
				"displayName": fixture.RandStringRunes(55),
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "Username Required",
			body: gin.H{
				"password":    user.Password,
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "DisplayName Required",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
				"email":    user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "Email Required",
			body: gin.H{
				"username":    user.Username,
				"password":    user.Password,
				"displayName": user.DisplayName,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
		{
			name: "Password Required",
			body: gin.H{
				"username":    user.Username,
				"displayName": user.DisplayName,
				"email":       user.Email,
			},
			buildStubs: func(mockUserService *mocks.UserService) {
				mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, mockUserService *mocks.UserService) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				mockUserService.AssertNotCalled(t, "Register")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			mockUserService := new(mocks.UserService)
			tc.buildStubs(mockUserService)

			// a response recorder for getting written http response
			rr := httptest.NewRecorder()

			router := gin.Default()
			store := cookie.NewStore([]byte("secret"))
			router.Use(sessions.Sessions("mqk", store))

			NewHandler(&Config{
				R:           router,
				UserService: mockUserService,
			})

			// create a request body with empty email and password
			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			// use bytes.NewBuffer to create a reader
			request, err := http.NewRequest(http.MethodPost, "/v1/accounts/register", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			tc.checkResponse(rr, mockUserService)
		})
	}
}
