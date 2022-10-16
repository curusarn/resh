package recio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/record"
	"github.com/curusarn/resh/internal/recordint"
)

// TODO: better errors
func (r *RecIO) WriteFile(fpath string, data []record.V1) error {
	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, rec := range data {
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

func (r *RecIO) EditRecordFlagsInFile(fpath string, idx int, rec recordint.Flag) error {
	// FIXME: implement
	// open file "not as append"
	// scan to the correct line

	return nil
}

func encodeV1Record(rec record.V1) ([]byte, error) {
	jsn, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to encode json: %w", err)
	}
	return append(jsn, []byte("\n")...), nil
}
