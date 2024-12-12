package middleware

import "github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"

type MiddlewareStruct struct {
	Logger  logger.LogrusLogger
	hashKey *string
}

func NewMiddlewareStruct(log logger.LogrusLogger, key *string) MiddlewareStruct {
	return MiddlewareStruct{
		Logger:  log,
		hashKey: key,
	}
}
