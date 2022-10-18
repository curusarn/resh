package recio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/curusarn/resh/internal/recconv"
	"github.com/curusarn/resh/internal/record"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
)

func (r *RecIO) ReadAndFixFile(fpath string, maxErrors int) ([]recordint.Indexed, error) {
	recs, numErrs, err := r.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	if numErrs > maxErrors {
		return nil, fmt.Errorf("encountered too many decoding errors")
	}
	if numErrs == 0 {
		return recs, nil
	}

	// TODO: check there error messages
	r.sugar.Warnw("Some history records could not be decoded - fixing resh history file by dropping them",
		"corruptedRecords", numErrs,
	)
	fpathBak := fpath + ".bak"
	r.sugar.Infow("Backing up current corrupted history file",
		"backupFilename", fpathBak,
	)
	// TODO: maybe use upstram copy function
	err = copyFile(fpath, fpathBak)
	if err != nil {
		r.sugar.Errorw("Failed to create a backup history file - aborting fixing history file",
			"backupFilename", fpathBak,
			zap.Error(err),
		)
		return recs, nil
	}
	r.sugar.Info("Writing resh history file without errors ...")
	var recsV1 []record.V1
	for _, rec := range recs {
		recsV1 = append(recsV1, rec.Rec)
	}
	err = r.OverwriteFile(fpath, recsV1)
	if err != nil {
		r.sugar.Errorw("Failed write fixed history file - aborting fixing history file",
			"filename", fpath,
			zap.Error(err),
		)
	}
	return recs, nil
}

func (r *RecIO) ReadFile(fpath string) ([]recordint.Indexed, int, error) {
	var recs []recordint.Indexed
	file, err := os.Open(fpath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	numErrs := 0
	var idx int
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		idx++
		rec, err := r.decodeLine(line)
		if err != nil {
			numErrs++
			continue
		}
		recidx := recordint.Indexed{
			Rec: *rec,
			// TODO: Is line index actually enough?
			// 	 Don't we want to count bytes because we will scan by number of bytes?
			// 	 hint: https://benjamincongdon.me/blog/2018/04/10/Counting-Scanned-Bytes-in-Go/
			Idx: idx,
		}
		recs = append(recs, recidx)
	}
	if err != io.EOF {
		r.sugar.Error("Error while loading file", zap.Error(err))
	}
	r.sugar.Infow("Loaded resh history records",
		"recordCount", len(recs),
	)
	return recs, numErrs, nil
}

func copyFile(source, dest string) error {
	from, err := os.Open(source)
	if err != nil {
		return err
	}
	defer from.Close()

	// This is equivalnet to: os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	to, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecIO) decodeLine(line string) (*record.V1, error) {
	idx := strings.Index(line, "{")
	if idx == -1 {
		return nil, fmt.Errorf("no openning brace found")
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
