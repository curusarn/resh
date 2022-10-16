package logger

import (
	"fmt"
	"path/filepath"

	"github.com/curusarn/resh/internal/datadir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(executable string, level zapcore.Level, developement bool) (*zap.Logger, error) {
	dataDir, err := datadir.GetPath()
	if err != nil {
		return nil, fmt.Errorf("error while getting resh data dir: %w", err)
	}
	logPath := filepath.Join(dataDir, "log.json")
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
