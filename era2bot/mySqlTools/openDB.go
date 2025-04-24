package mySqlTools

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	era2bot "github.com/disgoorg/bot-template/era2bot"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	configPath := "config.toml"
	cfg, err := era2bot.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))

	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.MySQL.DbUser, cfg.MySQL.DbPassword, cfg.MySQL.DbServer, cfg.MySQL.DbName)

	d, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))
	}

	// Test the connection
	if err := d.Ping(); err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))
	}

	// Optional: Set connection pool limits
	d.SetMaxOpenConns(25)
	d.SetMaxIdleConns(25)
	d.SetConnMaxLifetime(5 * time.Minute)
	DB = d
}
