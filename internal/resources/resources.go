package resources

import (
	"embed"
)

//go:embed config/app.*.yaml
var AppConfigFS embed.FS

//go:embed config/logger.*.yaml
var ZerologConfigFS embed.FS
