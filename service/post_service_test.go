package service

import (
	"fmt"
	"github.com/sentrionic/mirage/mocks"
	"github.com/sentrionic/mirage/model"
	"github.com/sentrionic/mirage/model/apperrors"
	"github.com/sentrionic/mirage/model/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPostService_FindPostByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()
		mockPost := fixture.GetMockPost()
		mockPost.ID = uid

		mockPostRepository := new(mocks.PostRepository)
		ps := NewPostService(&PSConfig{
			PostRepository: mockPostRepository,
		})
		mockPostRepository.On("FindByID", uid).Return(mockPost, nil)

		post, err := ps.FindPostByID(uid)

		assert.NoError(t, err)
		assert.Equal(t, post, mockPost)
		mockPostRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := GenerateId()

		mockPostRepository := new(mocks.PostRepository)
		ps := NewPostService(&PSConfig{
			PostRepository: mockPostRepository,
		})

		mockPostRepository.On("FindByID", uid).Return(nil, fmt.Errorf("some error down the call chain"))

		post, err := ps.FindPostByID(uid)

		assert.Nil(t, post)
		assert.Error(t, err)
		mockPostRepository.AssertExpectations(t)
	})
}

func TestPostService_CreatePost(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()
		mockPost := fixture.GetMockPost()

		initial := &model.Post{
			UserID: mockPost.UserID,
			Text:   mockPost.Text,
		}

		mockPostRepository := new(mocks.PostRepository)
		ps := NewPostService(&PSConfig{
			PostRepository: mockPostRepository,
		})

		mockPostRepository.
			On("Create", initial).
			Run(func(args mock.Arguments) {
				mockPost.ID = uid
			}).Return(mockPost, nil)

		post, err := ps.CreatePost(initial)

		assert.NoError(t, err)

		assert.Equal(t, uid, mockPost.ID)
		assert.Equal(t, post, mockPost)

		mockPostRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockPost := fixture.GetMockPost()
		initial := &model.Post{
			Text:   mockPost.Text,
			UserID: mockPost.UserID,
		}

		mockPostRepository := new(mocks.PostRepository)
		us := NewPostService(&PSConfig{
			PostRepository: mockPostRepository,
		})

		mockErr := apperrors.NewInternal()

		mockPostRepository.
			On("Create", initial).
			Return(nil, mockErr)

		post, err := us.CreatePost(initial)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())
		assert.Nil(t, post)

		mockPostRepository.AssertExpectations(t)
	})
}

func TestPostService_UploadFile(t *testing.T) {
	mockPostRepository := new(mocks.PostRepository)
	mockFileRepository := new(mocks.FileRepository)

	ps := NewPostService(&PSConfig{
		PostRepository: mockPostRepository,
		FileRepository: mockFileRepository,
	})

	mockUser := fixture.GetMockUser()

	t.Run("Successful image upload", func(t *testing.T) {
		mockPost := fixture.GetMockPost()

		file := &model.File{
			PostId:   mockPost.ID,
			FileType: "image/png",
		}

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "media/"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
			mock.AnythingOfType("string"),
			file.FileType,
		}

		imageURL := "https://imageurl.com/jdfkj34kljl"

		mockFileRepository.
			On("UploadFile", uploadFileArgs...).
			Return(imageURL, nil)

		file.Url = imageURL

		initial := &model.Post{
			File: file,
			User: *mockUser,
		}

		mockPostRepository.
			On("Create", initial).
			Return(mockPost, nil)

		uploadedFile, err := ps.UploadFile(imageFileHeader)
		assert.NoError(t, err)
		assert.NotNil(t, uploadedFile)

		newPost, err := ps.CreatePost(initial)

		assert.NoError(t, err)
		assert.Equal(t, newPost, mockPost)
		mockFileRepository.AssertCalled(t, "UploadFile", uploadFileArgs...)
		mockPostRepository.AssertCalled(t, "Create", initial)
	})

	t.Run("FileRepository Error", func(t *testing.T) {
		mockPostRepository := new(mocks.PostRepository)
		mockFileRepository := new(mocks.FileRepository)

		ps := NewPostService(&PSConfig{
			PostRepository: mockPostRepository,
			FileRepository: mockFileRepository,
		})

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "media/"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		}

		mockError := apperrors.NewInternal()
		mockFileRepository.
			On("UploadFile", uploadFileArgs...).
			Return("", mockError)

		uploadedFile, err := ps.UploadFile(imageFileHeader)
		assert.Nil(t, uploadedFile)
		assert.Error(t, err)

		mockFileRepository.AssertCalled(t, "UploadFile", uploadFileArgs...)
		mockPostRepository.AssertNotCalled(t, "Create")
	})

	t.Run("PostRepository Create Error", func(t *testing.T) {
		uid, _ := GenerateId()
		imageURL := "https://imageurl.com/jdfkj34kljl"

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "media/"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		}

		mockFileRepository.
			On("UploadFile", uploadFileArgs...).
			Return(imageURL, nil)

		file := &model.File{
			PostId:   uid,
			FileType: "image/png",
		}

		initial := &model.Post{
			ID:   uid,
			User: *mockUser,
			File: file,
		}

		mockError := apperrors.NewInternal()
		mockPostRepository.
			On("Create", initial).
			Return(nil, mockError)

		uploadedFile, err := ps.UploadFile(imageFileHeader)
		assert.NoError(t, err)
		assert.NotNil(t, uploadedFile)

		createdPost, err := ps.CreatePost(initial)

		assert.Error(t, err)
		assert.Nil(t, createdPost)
		mockFileRepository.AssertCalled(t, "UploadFile", uploadFileArgs...)
		mockPostRepository.AssertCalled(t, "Create", initial)
	})
}
