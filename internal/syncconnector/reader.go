package syncconnector

import "github.com/curusarn/resh/internal/record"

func (sc SyncConnector) getLatestRecord(machineId *string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (sc SyncConnector) downloadRecords(lastRecords map[string]string) ([]record.V1, error) {
	var records []record.V1
	return records, nil
}
