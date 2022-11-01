package histio

import (
	"path"

	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/record"
	"go.uber.org/zap"
)

type Histio struct {
	sugar   *zap.SugaredLogger
	histDir string

	thisDeviceID string
	thisHistory  *histfile
	// TODO: remote histories
	// moreHistories map[string]*histfile

	recordsToAppend chan record.V1
	recordsToFlag   chan recordint.Flag
}

func New(sugar *zap.SugaredLogger, dataDir, deviceID string) *Histio {
	sugarHistio := sugar.With(zap.String("component", "histio"))
	histDir := path.Join(dataDir, "history")
	currPath := path.Join(histDir, deviceID)
	// TODO: file extenstion for the history, yes or no? (<id>.reshjson vs. <id>)

	// TODO: discover other history files, exclude current

	return &Histio{
		sugar:   sugarHistio,
		histDir: histDir,

		thisDeviceID: deviceID,
		thisHistory:  newHistfile(sugar, currPath),
		// moreHistories: ...
	}
}

func (h *Histio) Append(r *record.V1) {

}
