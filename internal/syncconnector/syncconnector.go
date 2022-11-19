package syncconnector

import (
	"github.com/curusarn/resh/internal/histcli"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
	"net/url"
	"path"
	"time"
)

const storeEndpoint = "/store"
const historyEndpoint = "/history"
const latestEndpoint = "/latest"

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
			sc.sugar.Debug("checking remote")

			recs, err := sc.downloadRecords(sc.history.LatestRecordsPerDevice())
			if err != nil {
				continue
			}

			sc.sugar.Debugf("Got %d records", len(recs))

			for _, rec := range recs {
				sc.history.AddRecord(&recordint.Indexed{
					Rec: rec,
				})
			}

		}
	}(sc)

	return sc, nil
}

func (sc SyncConnector) getAddressWithPath(endpoint string) string {
	address := *sc.address
	address.Path = path.Join(address.Path, endpoint)
	return address.String()
}
