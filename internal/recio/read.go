package recio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/curusarn/resh/internal/futil"
	"github.com/curusarn/resh/internal/recconv"
	"github.com/curusarn/resh/record"
	"go.uber.org/zap"
)

func (r *RecIO) ReadAndFixFile(fpath string, maxErrors int) ([]record.V1, error) {
	recs, decodeErrs, err := r.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	numErrs := len(decodeErrs)
	if numErrs > maxErrors {
		r.sugar.Errorw("Encountered too many decoding errors",
			"errorsCount", numErrs,
			"individualErrors", "<Search 'Error while decoding line' to see individual errors>",
		)
		return nil, fmt.Errorf("encountered too many decoding errors, last error: %w", decodeErrs[len(decodeErrs)-1])
	}
	if numErrs == 0 {
		return recs, nil
	}

	r.sugar.Warnw("Some history records could not be decoded - fixing RESH history file by dropping them",
		"corruptedRecords", numErrs,
		"lastError", decodeErrs[len(decodeErrs)-1],
		"individualErrors", "<Search 'Error while decoding line' to see individual errors>",
	)

	fpathBak := fpath + ".bak"
	r.sugar.Infow("Backing up current corrupted history file",
		"historyFileBackup", fpathBak,
	)
	err = futil.CopyFile(fpath, fpathBak)
	if err != nil {
		r.sugar.Errorw("Failed to create a backup history file - aborting fixing history file",
			"historyFileBackup", fpathBak,
			zap.Error(err),
		)
		return recs, nil
	}
	r.sugar.Info("Writing resh history file without errors ...")
	err = r.OverwriteFile(fpath, recs)
	if err != nil {
		r.sugar.Errorw("Failed write fixed history file - restoring history file from backup",
			"historyFile", fpath,
			zap.Error(err),
		)

		err = futil.CopyFile(fpathBak, fpath)
		if err != nil {
			r.sugar.Errorw("Failed restore history file from backup",
				"historyFile", fpath,
				"HistoryFileBackup", fpathBak,
				zap.Error(err),
			)
		}
	}
	return recs, nil
}

func (r *RecIO) ReadFile(fpath string) ([]record.V1, []error, error) {
	var recs []record.V1
	file, err := os.Open(fpath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var decodeErrs []error
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		rec, err := r.decodeLine(line)
		if err != nil {
			r.sugar.Errorw("Error while decoding line", zap.Error(err),
				"filePath", fpath,
				"line", line,
			)
			decodeErrs = append(decodeErrs, err)
			continue
		}
		recs = append(recs, *rec)
	}
	if err != io.EOF {
		r.sugar.Error("Error while reading file", zap.Error(err))
		return recs, decodeErrs, err
	}
	r.sugar.Infow("Loaded resh history records",
		"recordCount", len(recs),
	)
	return recs, decodeErrs, nil
}

func (r *RecIO) decodeLine(line string) (*record.V1, error) {
	idx := strings.Index(line, "{")
	if idx == -1 {
		return nil, fmt.Errorf("no opening brace found")
	}
	schema := line[:idx]
	jsn := line[idx:]
	switch schema {
	case "v1":
		var rec record.V1
		err := decodeAnyRecord(jsn, &rec)
		if err != nil {
			return nil, err
		}
		return &rec, nil
	case "":
		var rec record.Legacy
		err := decodeAnyRecord(jsn, &rec)
		if err != nil {
			return nil, err
		}
		return recconv.LegacyToV1(&rec), nil
	default:
		return nil, fmt.Errorf("unknown record schema/type '%s'", schema)
	}
}

// TODO: find out if we are loosing performance because of the use of interface{}

func decodeAnyRecord(jsn string, rec interface{}) error {
	err := json.Unmarshal([]byte(jsn), &rec)
	if err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}
	return nil
}
