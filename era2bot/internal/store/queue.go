package store

import (
	"database/sql"
	"fmt"

	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/mySqlTools"
)

type QueueStore struct {
	DB *sql.DB
}

// SELECT global queue
func (s *QueueStore) GetQueue(curator *models.Curator) (*models.Queue, error) {
	queue := &models.Queue{}
	queue.Type = models.GlobalQueue
	queueFilter := "%"
	if curator != nil {
		queueFilter = curator.Name
		queue.Type = models.PersonalQueue
	}
	fmt.Println("store query")
	err := mySqlTools.DB.QueryRow("SELECT count(*) from curationChoices WHERE voted = 0 and curator like ?", queueFilter).
		Scan(&queue.Count)
	if err != nil {
		fmt.Println("store query error")
		return nil, err
	}
	fmt.Println("store query done")
	return queue, nil
}

// // UPDATE
// func (s *QueueStore) UpdateEmail(id int, newEmail string) error {
// 	_, err := s.DB.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, id)
// 	return err
// }
