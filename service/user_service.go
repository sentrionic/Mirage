package service

import (
	"github.com/sentrionic/mirage/model"
	"github.com/sentrionic/mirage/model/apperrors"
	"log"
	"mime/multipart"
)

type userService struct {
	UserRepository model.UserRepository
	FileRepository model.FileRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository model.UserRepository
	FileRepository model.FileRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
		FileRepository: c.FileRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *userService) Get(uid string) (*model.User, error) {
	user, err := s.UserRepository.FindByID(uid)

	return user, err
}

func (s *userService) Register(user *model.User) (*model.User, error) {
	pw, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}

	user.Password = pw

	id, err := GenerateId()

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}
	user.ID = id

	user.Image = GetGravatar(user.Email)

	return s.UserRepository.Create(user)
}

func (s *userService) Login(email, password string) (*model.User, error) {
	user, err := s.UserRepository.FindByEmail(email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return nil, apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method
	match, err := comparePasswords(user.Password, password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization("Invalid email and password combination")
	}

	return user, nil
}

func (s *userService) Update(user *model.User) error {
	return s.UserRepository.Update(user)
}

func (s *userService) ChangeAvatar(header *multipart.FileHeader, directory string) (string, error) {
	return s.FileRepository.UploadAvatar(header, directory)
}

func (s *userService) DeleteImage(key string) error {
	return s.FileRepository.DeleteImage(key)
}
