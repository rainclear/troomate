package config

import (
	"database/sql"
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	InfoLog      *log.Logger
	InProduction bool
	DBPath       string
	Db           *sql.DB
	Accounts     []string
}
