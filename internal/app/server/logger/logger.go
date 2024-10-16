package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var once sync.Once

func InstanceLogger(level string) *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()

		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			logger.Errorf("can not parse log level: %v; current log level is %s", err, logger.GetLevel().String())
			return
		}

		logger.SetLevel(lvl)
	})

	return logger
}
