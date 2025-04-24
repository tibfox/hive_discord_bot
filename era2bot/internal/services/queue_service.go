package services

import (
	"fmt"

	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/internal/store"
)

type QueueService struct {
	Store *store.QueueStore
}

// Factory function: builds the store inside
func NewQueueService() *QueueService {
	return &QueueService{
		Store: &store.QueueStore{},
	}
}

func (s *QueueService) GetQueue(queueType models.QueueType, curator *models.Curator) (*models.Queue, error) {

	queue, err := s.Store.GetQueue(curator)
	if err != nil {
		return nil, err

	} else {
		if queue == nil {
			return nil, fmt.Errorf("failed to get queue")
		} else {
			return queue, nil
		}

	}
}

// func (s *UserService) ChangeEmail(id int, email string) error {
//     return s.Users.UpdateEmail(id, email)
// }
