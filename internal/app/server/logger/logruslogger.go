package logger

import (
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	LogrusLog *logrus.Logger
}

func NewLogrusLogger() LogrusLogger {
	return LogrusLogger{
		LogrusLog: logrus.New(),
	}
}

func (llog LogrusLogger) SetLevelForLog(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	llog.LogrusLog.SetLevel(lvl)
	return nil
}
