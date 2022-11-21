package syncconnector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/curusarn/resh/record"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (sc SyncConnector) getLatestRecord(machineId *string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (sc SyncConnector) downloadRecords(lastRecords map[string]float64) ([]record.V1, error) {
	var records []record.V1

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	latestRes := map[string]string{}
	for device, t := range lastRecords {
		sc.sugar.Debugf("Latest for %s is %f", device, t)
		latestRes[device] = fmt.Sprintf("%.4f", t)
	}

	latestJson, err := json.Marshal(latestRes)
	if err != nil {
		sc.sugar.Errorw("converting latest to JSON failed", "err", err)
		return nil, err
	}
	reqBody := bytes.NewBuffer(latestJson)

	address := sc.getAddressWithPath(historyEndpoint)
	resp, err := client.Post(address, "application/json", reqBody)
	if err != nil {
		sc.sugar.Errorw("history request failed", "address", address, "err", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sc.sugar.Errorw("reader close failed", "err", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sc.sugar.Warnw("reading response body failed", "err", err)
	}

	err = json.Unmarshal(body, &records)
	if err != nil {
		sc.sugar.Errorw("Unmarshalling failed", "err", err)
		return nil, err
	}

	return records, nil
}

func (sc SyncConnector) latest() (map[string]float64, error) {
	var knownDevices []string
	for deviceId, _ := range sc.history.LatestRecordsPerDevice() {
		knownDevices = append(knownDevices, deviceId)
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	knownJson, err := json.Marshal(knownDevices)
	if err != nil {
		sc.sugar.Errorw("converting latest to JSON failed", "err", err)
		return nil, err
	}
	reqBody := bytes.NewBuffer(knownJson)

	address := sc.getAddressWithPath(latestEndpoint)
	resp, err := client.Post(address, "application/json", reqBody)
	if err != nil {
		sc.sugar.Errorw("latest request failed", "address", address, "err", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sc.sugar.Errorw("reader close failed", "err", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sc.sugar.Warnw("reading response body failed", "err", err)
	}

	latest := map[string]string{}

	err = json.Unmarshal(body, &latest)
	if err != nil {
		sc.sugar.Errorw("Unmarshalling failed", "err", err)
		return nil, err
	}

	l := make(map[string]float64, len(latest))
	for deviceId, ts := range latest {
		t, err := strconv.ParseFloat(ts, 64)
		if err != nil {
			return nil, err
		}
		l[deviceId] = t
	}

	return l, nil
}
