package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/mySqlTools"
)

type AuthorStore struct {
	DB *sql.DB
}

// SELECT Author by name
func (s *AuthorStore) GetAuthor(authorName *string) (*models.Author, error) {
	author := &models.Author{
		Name: *authorName,
	}
	var reason string = ""
	var curator string = ""
	var blacklistedOn time.Time

	err := mySqlTools.DB.QueryRow("SELECT reason, curator, blacklistedOn from userBlacklist where username = ?", authorName).
		Scan(&reason, &curator, &blacklistedOn)

	if err != nil {
		if err == sql.ErrNoRows {
			return author, nil
		} else {
			fmt.Println("store query error")
			return nil, err
		}

	} else {
		if curator != "" {
			author.BlacklistedBy = curator
			author.BlacklistedOn = blacklistedOn
			author.BlacklistReason = reason
		}
	}
	fmt.Println("store query done")
	return author, nil
}
