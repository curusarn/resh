package syncconnector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/curusarn/resh/record"
	"io"
	"log"
	"net/http"
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
	responseBody := bytes.NewBuffer(latestJson)

	address := sc.getAddressWithPath(historyEndpoint)
	resp, err := client.Post(address, "application/json", responseBody)
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
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &records)
	if err != nil {
		sc.sugar.Errorw("Unmarshalling failed", "err", err)
		return nil, err
	}

	return records, nil
}

func latest() {
	//curl localhost:8080/latest -X POST -d '[]'
	//curl localhost:8080/latest -X POST -d '["one"]'
}
