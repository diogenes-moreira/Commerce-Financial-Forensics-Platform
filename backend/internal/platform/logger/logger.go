package logger

import (
	"github.com/diogenes/costforensics/backend/config"
	"github.com/sirupsen/logrus"
)

func Setup(cfg config.LogConfig) *logrus.Logger {
	log := logrus.New()

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log
}
