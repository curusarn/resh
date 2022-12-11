package recio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/record"
)

// TODO: better errors
func (r *RecIO) OverwriteFile(fpath string, recs []record.V1) error {
	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeRecords(file, recs)
}

// TODO: better errors
func (r *RecIO) AppendToFile(fpath string, recs []record.V1) error {
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeRecords(file, recs)
}

// TODO: better errors
// TODO: rethink this
func (r *RecIO) EditRecordFlagsInFile(fpath string, idx int, rec recordint.Flag) error {
	// FIXME: implement
	// open file "not as append"
	// scan to the correct line
	r.sugar.Error("not implemented yet (FIXME)")
	return nil
}

func writeRecords(file *os.File, recs []record.V1) error {
	for _, rec := range recs {
		jsn, err := encodeV1Record(rec)
		if err != nil {
			return err
		}
		_, err = file.Write(jsn)
		if err != nil {
			return err
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
