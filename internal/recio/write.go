package recio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/curusarn/resh/record"
)

func (r *RecIO) OverwriteFile(fpath string, recs []record.V1) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("could not create/truncate file: %w", err)
	}
	err = writeRecords(file, recs)
	if err != nil {
		return fmt.Errorf("error while writing records: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}
	return nil
}

func (r *RecIO) AppendToFile(fpath string, recs []record.V1) error {
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open/create file: %w", err)
	}
	err = writeRecords(file, recs)
	if err != nil {
		return fmt.Errorf("error while writing records: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}
	return nil
}

func writeRecords(file *os.File, recs []record.V1) error {
	for _, rec := range recs {
		jsn, err := encodeV1Record(rec)
		if err != nil {
			return fmt.Errorf("could not encode record: %w", err)
		}
		_, err = file.Write(jsn)
		if err != nil {
			return fmt.Errorf("could not write json: %w", err)
		}
	}
	return nil
}

func encodeV1Record(rec record.V1) ([]byte, error) {
	version := []byte("v1")
	jsn, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to encode json: %w", err)
	}
	return append(append(version, jsn...), []byte("\n")...), nil
}
