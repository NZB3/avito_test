package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string)
	Infof(format string, args ...any)
	Error(msg string)
	Errorf(format string, args ...any)
}

type logger struct {
	l *slog.Logger
}

func New() *logger {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &logger{
		l: l,
	}
}

func (log *logger) Info(msg string) {
	log.l.Info(msg)
}

func (log *logger) Infof(format string, args ...any) {
	log.l.Info(fmt.Sprintf(format, args...))
}

func (log *logger) Error(msg string) {
	log.l.Error(msg)
}

func (log *logger) Errorf(format string, args ...any) {
	log.l.Error(fmt.Sprintf(format, args...))
}
