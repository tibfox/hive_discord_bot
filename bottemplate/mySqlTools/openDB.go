package mySqlTools

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/disgoorg/bot-template/bottemplate"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() bool {
	configPath := "config.toml"
	cfg, err := bottemplate.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))
		return false
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.MySQL.DbUser, cfg.MySQL.DbPassword, cfg.MySQL.DbServer, cfg.MySQL.DbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("error opening db %s", err)
		return false
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(); err != nil {
		fmt.Println("error pinging db %s", err)
		return false

	}
	return true
}

// rows, err := db.Query("SELECT id, title, created_at FROM posts")
// if err != nil {
//     panic(err)
// }
// defer rows.Close()

// for rows.Next() {
//     var id int
//     var title string
//     var createdAt string // or time.Time if you prefer

//     if err := rows.Scan(&id, &title, &createdAt); err != nil {
//         panic(err)
//     }

//     fmt.Printf("ID: %d, Title: %s, Created At: %s\n", id, title, createdAt)
// }
