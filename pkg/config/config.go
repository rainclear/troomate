package config

import (
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	InfoLog      *log.Logger
	InProduction bool
}
