package service

import (
	"fmt"
	"github.com/lucsky/cuid"
	"github.com/sentrionic/mirage/model"
	"github.com/sentrionic/mirage/model/apperrors"
	"log"
	"mime/multipart"
	"path"
)

type postService struct {
	PostRepository model.PostRepository
	FileRepository model.FileRepository
}

// PSConfig will hold repositories that will eventually be injected into this
// this service layer
type PSConfig struct {
	PostRepository model.PostRepository
	FileRepository model.FileRepository
}

// NewPostService is a factory function for
// initializing a PostService with its repository layer dependencies
func NewPostService(c *PSConfig) model.PostService {
	return &postService{
		PostRepository: c.PostRepository,
		FileRepository: c.FileRepository,
	}
}

func (p *postService) FindPostByID(id string) (*model.Post, error) {
	return p.PostRepository.FindByID(id)
}

func (p *postService) CreatePost(post *model.Post) (*model.Post, error) {
	id, err := GenerateId()

	if err != nil {
		log.Printf("Unable to create post for author: %v\n", post.UserID)
		return nil, apperrors.NewInternal()
	}
	post.ID = id

	if post.Text != nil {
		post.HashTags = GetHashtags(*post.Text)
	}

	return p.PostRepository.Create(post)
}

func (p *postService) DeletePost(post *model.Post) error {
	return p.PostRepository.Delete(post)
}

func (p *postService) UploadFile(header *multipart.FileHeader) (*model.File, error) {
	slug := cuid.New()
	ext := path.Ext(header.Filename)
	filename := slug + ext
	mimetype := header.Header.Get("Content-Type")

	file := model.File{
		FileType: mimetype,
		Filename: filename,
	}

	id, err := GenerateId()
	if err != nil {
		return nil, err
	}

	file.ID = id

	directory := fmt.Sprintf("media/")
	url, err := p.FileRepository.UploadFile(header, directory, filename, mimetype)

	if err != nil {
		return nil, err
	}

	file.Url = url

	return &file, nil
}
