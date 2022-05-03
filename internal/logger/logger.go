package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(executable string, level zapcore.Level, developement bool) (*zap.Logger, error) {
	// TODO: consider getting log path from config ?
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting home dir: %w", err)
	}
	logPath := filepath.Join(homeDir, ".resh/log.json")
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{logPath}
	loggerConfig.Level.SetLevel(level)
	loggerConfig.Development = developement // DPanic panics in developement
	logger, err := loggerConfig.Build()
	if err != nil {
		return logger, fmt.Errorf("error while creating logger: %w", err)
	}
	return logger.With(zap.String("executable", executable)), err
}
