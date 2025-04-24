package services

import (
	"fmt"

	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/internal/store"
)

type AuthorService struct {
	Store *store.AuthorStore
}

// Factory function: builds the store inside
func NewAuthorService() *AuthorService {
	return &AuthorService{
		Store: &store.AuthorStore{},
	}
}

func (s *AuthorService) GetAuthor(authorName *string) (*models.Author, error) {

	Author, err := s.Store.GetAuthor(authorName)
	if err != nil {
		return nil, err

	} else {
		if Author == nil {
			return nil, fmt.Errorf("failed to get Author")
		} else {
			return Author, nil
		}

	}
}

// func (s *UserService) ChangeEmail(id int, email string) error {
//     return s.Users.UpdateEmail(id, email)
// }
