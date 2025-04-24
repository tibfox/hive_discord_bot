package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/mySqlTools"
)

type CuratorStore struct {
	DB *sql.DB
}

// SELECT curator by name
func (s *CuratorStore) GetCuratorByName(curatorName *string) (*models.Curator, error) {
	curator := &models.Curator{
		Name: *curatorName,
	}
	yesterday := time.Now().Add(-24 * time.Hour)

	err := mySqlTools.DB.QueryRow(`
		SELECT discorduser, (
			select 	count(*)
			from	curationChoices
			where	STR_TO_DATE(CONCAT(curationDate, ' ', curationTime), '%Y-%m-%d %H:%i:%s') > ?
			and curator = curatorProfiles.curator
			)
		from curatorProfiles
		where curator = ?
		`, yesterday, curatorName).
		Scan(&curator.DiscordUserId, &curator.CountCurationsInInterval)
	if err != nil {
		fmt.Println("store query error")
		return nil, err
	}
	fmt.Println("store query done")
	return curator, nil
}

// SELECT curator by name
func (s *CuratorStore) GetCuratorByDiscordId(discordId *string) (*models.Curator, error) {
	curator := &models.Curator{
		DiscordUserId: *discordId,
	}

	yesterday := time.Now().Add(-24 * time.Hour)

	err := mySqlTools.DB.QueryRow(`
	SELECT curator, (
			select 	count(*)
			from	curationChoices
			where	STR_TO_DATE(CONCAT(curationDate, ' ', curationTime), '%Y-%m-%d %H:%i:%s') > ?
			and curator = curatorProfiles.curator
			)
	from curatorProfiles 
	where discorduser = ?`, yesterday, discordId).
		Scan(&curator.Name, &curator.CountCurationsInInterval)
	if err != nil {
		fmt.Println("store query error %w", err)
		return nil, err
	}
	fmt.Println("store query done", curator)
	return curator, nil
}

// // UPDATE
// func (s *QueueStore) UpdateEmail(id int, newEmail string) error {
// 	_, err := s.DB.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, id)
// 	return err
// }
