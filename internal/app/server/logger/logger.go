package logger

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
}

func SetLevelForLog(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	logger.SetLevel(lvl)
	return nil
}

func Logger() *logrus.Logger {
	return logger
}
