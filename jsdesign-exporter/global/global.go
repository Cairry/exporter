package global

import (
	"exporter/config"
	"go.uber.org/zap"
)

var (
	GvaServerConfig *config.ServerConfig
	GvaLogger       *zap.Logger
)
