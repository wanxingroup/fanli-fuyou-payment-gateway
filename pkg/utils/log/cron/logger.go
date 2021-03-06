package cron

/**
 * 使用 Logrus 实现了支持cron的日志接口
 */

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Entry
}

func NewLogger(logger *logrus.Entry) *Logger {

	return &Logger{logger: logger}
}

// Info logs routine messages about cron's operation.
func (l *Logger) Info(message string, keysAndValues ...interface{}) {

	l.logger.WithField("parameters", keysAndValues).Info(message)
}

// Error logs an error condition.
func (l *Logger) Error(err error, message string, keysAndValues ...interface{}) {

	l.logger.WithError(err).WithField("parameters", keysAndValues).Error(message)
}
