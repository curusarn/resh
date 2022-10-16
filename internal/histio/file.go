package histio

import (
	"fmt"
	"os"
	"sync"

	"github.com/curusarn/resh/internal/recio"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
)

type histfile struct {
	sugar *zap.SugaredLogger
	// deviceID string
	path string

	mu       sync.RWMutex
	data     []recordint.Indexed
	fileinfo os.FileInfo
}

func newHistfile(sugar *zap.SugaredLogger, path string) *histfile {
	return &histfile{
		sugar: sugar.With(
			// FIXME: drop V1 once original histfile is gone
			"component", "histfileV1",
			"path", path,
		),
		// deviceID: deviceID,
		path: path,
	}
}

func (h *histfile) updateFromFile() error {
	rio := recio.New(h.sugar)
	// TODO: decide and handle errors
	newData, _, err := rio.ReadFile(h.path)
	if err != nil {
		return fmt.Errorf("could not read history file: %w", err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data = newData
	h.updateFileInfo()
	return nil
}

func (h *histfile) updateFileInfo() error {
	info, err := os.Stat(h.path)
	if err != nil {
		return fmt.Errorf("history file not found: %w", err)
	}
	h.fileinfo = info
	return nil
}
