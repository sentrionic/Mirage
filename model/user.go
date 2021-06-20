package model

import (
	"mime/multipart"
	"time"
)

type AccountResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	Image       string  `json:"image"`
	Bio         *string `json:"about"`
}

func (user *User) NewAccountResponse() AccountResponse {
	return AccountResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Image:       user.Image,
		Bio:         user.Bio,
	}
}

type User struct {
	ID          string `gorm:"primaryKey"`
	Username    string `gorm:"not null;index;uniqueIndex"`
	DisplayName string `gorm:"not null;index"`
	Email       string `gorm:"not null;uniqueIndex"`
	Password    string `gorm:"not null" json:"-"`
	Image       string `gorm:"not null"`
	Bio         *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserService interface {
	Get(uid string) (*User, error)
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
	Update(user *User) error
	ChangeAvatar(header *multipart.FileHeader, directory string) (string, error)
	DeleteImage(key string) error
}

type UserRepository interface {
	FindByID(uid string) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) (*User, error)
	Update(user *User) error
}
