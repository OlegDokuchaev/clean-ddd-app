package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
	Debug(msg string, fields map[string]any)
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
}

type LoggerImpl struct {
	log *logrus.Logger
}

func NewLogger(log *logrus.Logger) Logger {
	return &LoggerImpl{log: log}
}

func (l *LoggerImpl) Info(msg string, fields map[string]any) {
	entry := l.log.WithFields(logrus.Fields{})
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Info(msg)
}

func (l *LoggerImpl) Warn(msg string, fields map[string]any) {
	entry := l.log.WithFields(logrus.Fields{})
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Warn(msg)
}

func (l *LoggerImpl) Error(msg string, fields map[string]any) {
	entry := l.log.WithFields(logrus.Fields{})
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Error(msg)
}

func (l *LoggerImpl) Debug(msg string, fields map[string]any) {
	entry := l.log.WithFields(logrus.Fields{})
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Debug(msg)
}

func (l *LoggerImpl) Printf(format string, v ...any) {
	l.log.Infof(format, v...)
}

func (l *LoggerImpl) Println(v ...any) {
	l.log.Infoln(v...)
}

func (l *LoggerImpl) Fatalf(format string, v ...any) {
	l.log.Fatalf(format, v...)
}

var _ Logger = (*LoggerImpl)(nil)
