package discogpt

import "go.uber.org/zap"

// wrap zap in an interface for mocking
//
//go:generate mockgen -source ./log.go -destination ./mock/log.go
type Logger interface {
	Infof(template string, args ...any)
	Debugf(template string, args ...any)
	Errorf(template string, arg ...any)
	Sync() error
}

func NewLogger() (*zap.SugaredLogger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}
