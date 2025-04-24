package services

import (
	"fmt"

	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/internal/store"
)

type CuratorService struct {
	Store *store.CuratorStore
}

// Factory function: builds the store inside
func NewCuratorService() *CuratorService {
	return &CuratorService{
		Store: &store.CuratorStore{},
	}
}

func (s *CuratorService) GetCurator(name *string, discordUserId *string) (*models.Curator, error) {
	if name != nil {
		curator, err := s.Store.GetCuratorByName(name)
		if err != nil {
			return nil, err

		} else {
			return curator, nil
		}
	} else {
		if discordUserId != nil {
			curator, err := s.Store.GetCuratorByDiscordId(discordUserId)
			if err != nil {
				return nil, err
			} else {
				if curator == nil {
					return nil, fmt.Errorf("failed to get user for %v", discordUserId)
				} else {
					return curator, nil
				}
			}
		}
	}
	return nil, nil
}
