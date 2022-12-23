package syncconnector

import (
	"bytes"
	"encoding/json"
	"github.com/curusarn/resh/record"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (sc SyncConnector) write() error {
	latestRemote, err := sc.latest()
	if err != nil {
		return err
	}
	latestLocal := sc.history.LatestRecordsPerDevice()
	remoteIsOlder := false
	for deviceId, lastLocal := range latestLocal {
		if lastRemote, ok := latestRemote[deviceId]; !ok {
			// Unknown deviceId on the remote - add records have to be sent
			remoteIsOlder = true
			break
		} else if lastLocal > lastRemote {
			remoteIsOlder = true
			break
		}
	}
	if !remoteIsOlder {
		sc.sugar.Debug("No need to sync remote, there are no newer local records")
		return nil
	}
	var toSend []record.V1
	for _, r := range sc.history.DumpRaw() {
		t, err := strconv.ParseFloat(r.Time, 64)
		if err != nil {
			sc.sugar.Warnw("Invalid time for record - skipping", "time", r.Time)
			continue
		}
		l, ok := latestRemote[r.DeviceID]
		if ok && l >= t {
			continue
		}
		sc.sugar.Infow("record is newer", "new", t, "old", l, "id", r.RecordID, "deviceid", r.DeviceID)
		toSend = append(toSend, r)
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	toSendJson, err := json.Marshal(toSend)
	if err != nil {
		sc.sugar.Errorw("converting toSend to JSON failed", "err", err)
		return err
	}
	reqBody := bytes.NewBuffer(toSendJson)

	address := sc.getAddressWithPath(storeEndpoint)
	resp, err := client.Post(address, "application/json", reqBody)
	if err != nil {
		sc.sugar.Errorw("store request failed", "address", address, "err", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sc.sugar.Errorw("reader close failed", "err", err)
		}
	}(resp.Body)

	sc.sugar.Debugw("store call", "status", resp.Status)

	return nil
}
