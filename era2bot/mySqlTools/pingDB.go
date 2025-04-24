package mySqlTools

import "log/slog"

func PingDB() error {
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))
		return err
	} else {
		return nil
	}
}
