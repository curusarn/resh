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
}

func New(sugar *zap.SugaredLogger, address string, authToken string, pullPeriodSeconds int, sendPeriodSeconds int, history *histcli.Histcli) (*SyncConnector, error) {
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
			sc.sugar.Debug("checking remote for new records")

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

	go func(sc *SyncConnector) {
		// wait to properly load all the records
		time.Sleep(time.Second * time.Duration(sendPeriodSeconds))
		for _ = range time.Tick(time.Second * time.Duration(sendPeriodSeconds)) {
			sc.sugar.Debug("syncing local records to the remote")

			err := sc.write()
			if err != nil {
				sc.sugar.Warnw("sending records to the remote failed", "err", err)
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
