package model

import (
	"github.com/lib/pq"
	"mime/multipart"
	"time"
)

type PostResponse struct {
	ID     string  `json:"id"`
	Text   *string `json:"text"`
	File   *File   `json:"file"`
	Author Profile `json:"author"`
}

func (post *Post) NewPostResponse(id string) PostResponse {
	return PostResponse{
		ID:     post.ID,
		Text:   post.Text,
		File:   post.File,
		Author: post.User.NewProfileResponse(id),
	}
}

type Post struct {
	ID        string `gorm:"primaryKey"`
	Text      *string
	File      *File          `gorm:"constraint:OnDelete:CASCADE;"`
	HashTags  pq.StringArray `gorm:"type:text[]"`
	UserID    string         `gorm:"not null;constraint:OnDelete:CASCADE;"`
	User      User           `gorm:"not null;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time
}

type PostService interface {
	FindPostByID(id string) (*Post, error)
	CreatePost(post *Post) (*Post, error)
	DeletePost(post *Post) error
	UploadFile(header *multipart.FileHeader) (*File, error)
}

type PostRepository interface {
	FindByID(id string) (*Post, error)
	Create(post *Post) (*Post, error)
	Delete(post *Post) error
}
