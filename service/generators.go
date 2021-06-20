package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bwmarrin/snowflake"
)

// GenerateId generates a snowflake id
func GenerateId() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Generate a snowflake ID.
	id := node.Generate()

	return id.String(), nil
}

// GetGravatar returns a link to a gravatar for the given email
// e.g. https://gravatar.com/avatar/55502f40dc8b7c769880b10874abc9d0?d=identicon
func GetGravatar(email string) string {
	hash := md5.Sum([]byte(email))
	value := hex.EncodeToString(hash[:])
	return fmt.Sprintf("https://gravatar.com/avatar/%s?d=identicon", value)
}
