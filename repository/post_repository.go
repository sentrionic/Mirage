package repository

import (
	"errors"
	"github.com/sentrionic/mirage/model"
	"github.com/sentrionic/mirage/model/apperrors"
	"gorm.io/gorm"
	"log"
)

// postRepository is data/repository implementation
// of service layer PostRepository
type postRepository struct {
	DB *gorm.DB
}

// NewPostRepository is a factory for initializing Post Repositories
func NewPostRepository(db *gorm.DB) model.PostRepository {
	return &postRepository{
		DB: db,
	}
}

// FindByID returns a post for the given ID
func (r *postRepository) FindByID(id string) (*model.Post, error) {
	post := &model.Post{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, apperrors.NewNotFound("id", id)
		}
		return post, apperrors.NewInternal()
	}

	return post, nil
}

// Create inserts the post in the DB
func (r *postRepository) Create(post *model.Post) (*model.Post, error) {
	if result := r.DB.Create(&post); result.Error != nil {
		log.Printf("Could not create a post for author: %v. Reason: %v\n", post.UserID, result.Error)
		return nil, apperrors.NewInternal()
	}

	return post, nil
}

func (r *postRepository) Delete(post *model.Post) error {
	return r.DB.Delete(&post).Error
}
