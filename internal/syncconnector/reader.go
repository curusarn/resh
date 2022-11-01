package syncconnector

import (
	"bytes"
	"encoding/json"
	"github.com/curusarn/resh/internal/record"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (sc SyncConnector) getLatestRecord(machineId *string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (sc SyncConnector) downloadRecords(lastRecords map[string]string) ([]record.V1, error) {
	var records []record.V1

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	// TODO: create request based on the local last records
	responseBody := bytes.NewBuffer([]byte("{}"))

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
	body, err := ioutil.ReadAll(resp.Body)
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
