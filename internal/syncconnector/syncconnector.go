package syncconnector

import (
	"github.com/curusarn/resh/internal/histcli"
	"github.com/curusarn/resh/internal/record"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type SyncConnector struct {
	sugar *zap.SugaredLogger

	address   *url.URL
	authToken string

	history *histcli.Histcli

	// TODO periodic push (or from the write channel)
	// TODO push period
}

func New(sugar *zap.SugaredLogger, address string, authToken string, pullPeriodSeconds int, history *histcli.Histcli) (*SyncConnector, error) {
	parsedAddress, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	sc := &SyncConnector{
		sugar:     sugar.With(zap.String("component", "syncConnector")),
		authToken: authToken,
		address:   parsedAddress,
		history:   history,
	}

	// TODO: propagate signals
	go func(sc *SyncConnector) {
		for _ = range time.Tick(time.Second * time.Duration(pullPeriodSeconds)) {
			sc.sugar.Infow("checking remote")

			// Add fake record (this will be produced by the sync connector)
			sc.history.AddRecord(&recordint.Indexed{
				Rec: record.V1{
					CmdLine:  "__fake_test__",
					DeviceID: "__test__",
				},
			})

		}
	}(sc)

	return sc, nil
}
