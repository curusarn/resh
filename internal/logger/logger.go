package logger

import (
	"fmt"
	"path/filepath"

	"github.com/curusarn/resh/internal/datadir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(executable string, level zapcore.Level, development string) (*zap.Logger, error) {
	dataDir, err := datadir.MakePath()
	if err != nil {
		return nil, fmt.Errorf("error while getting RESH data dir: %w", err)
	}
	logPath := filepath.Join(dataDir, "log.json")
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{logPath}
	loggerConfig.Level.SetLevel(level)
	loggerConfig.Development = development == "true" // DPanic panics in development
	logger, err := loggerConfig.Build()
	if err != nil {
		return logger, fmt.Errorf("error while creating logger: %w", err)
	}
	return logger.With(zap.String("executable", executable)), err
}
