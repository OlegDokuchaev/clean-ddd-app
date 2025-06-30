package logger

import (
	"github.com/sirupsen/logrus"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

type Logger interface {
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
	Debug(msg string, fields map[string]any)
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Log(level Level, msg string, fields map[string]any)
}

type LoggerImpl struct {
	log *logrus.Logger
}

func NewLogger(log *logrus.Logger) Logger {
	return &LoggerImpl{log: log}
}

func (l *LoggerImpl) logWithFields(fields map[string]any) *logrus.Entry {
	if fields != nil {
		return l.log.WithFields(fields)
	}
	return logrus.NewEntry(l.log)
}

func (l *LoggerImpl) Info(msg string, fields map[string]any) {
	l.logWithFields(fields).Info(msg)
}

func (l *LoggerImpl) Warn(msg string, fields map[string]any) {
	l.logWithFields(fields).Warn(msg)
}

func (l *LoggerImpl) Error(msg string, fields map[string]any) {
	l.logWithFields(fields).Error(msg)
}

func (l *LoggerImpl) Debug(msg string, fields map[string]any) {
	l.logWithFields(fields).Debug(msg)
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

func (l *LoggerImpl) Log(level Level, msg string, fields map[string]any) {
	switch level {
	case Debug:
		l.Debug(msg, fields)
	case Info:
		l.Info(msg, fields)
	case Warn:
		l.Warn(msg, fields)
	case Error:
		l.Error(msg, fields)
	case Fatal:
		l.Fatalf("%s", msg)
	default:
		l.Info(msg, fields)
	}
}

var _ Logger = (*LoggerImpl)(nil)
