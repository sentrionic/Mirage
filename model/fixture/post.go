package fixture

import (
	"github.com/sentrionic/mirage/model"
	"time"
)

func GetMockPost() *model.Post {
	text := RandStringRunes(60)
	return &model.Post{
		ID:        RandID(),
		Text:      &text,
		UserID:    RandID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
