package config

import "time"

const (
	DEV = "development"
	PRD = "production"
)

const (
	APP_CONFIG_PATH     = "config/app.%s.yaml"
	ZEROLOG_CONFIG_PATH = "config/logger.%s.yaml"
)

const TIME_FORMAT = time.RFC3339
