package recordint

import (
	"github.com/curusarn/resh/internal/record"
	"github.com/curusarn/resh/internal/recutil"
)

// TODO: This all seems excessive
// TODO: V1 should be converted directly to SearchApp record

// EnrichedRecord - record enriched with additional data
type Enriched struct {
	// TODO: think about if it really makes sense to have this based on V1
	record.V1

	// TODO: drop some/all of this
	// enriching fields - added "later"
	Command             string   `json:"command"`
	FirstWord           string   `json:"firstWord"`
	Invalid             bool     `json:"invalid"`
	SeqSessionID        uint64   `json:"seqSessionId"`
	LastRecordOfSession bool     `json:"lastRecordOfSession"`
	DebugThisRecord     bool     `json:"debugThisRecord"`
	Errors              []string `json:"errors"`
	// SeqSessionID uint64 `json:"seqSessionId,omitempty"`
}

// Enriched - returns enriched record
func NewEnrichedFromV1(r *record.V1) Enriched {
	rec := Enriched{Record: r}
	// normlize git remote
	rec.GitOriginRemote = NormalizeGitRemote(rec.GitOriginRemote)
	rec.GitOriginRemoteAfter = NormalizeGitRemote(rec.GitOriginRemoteAfter)
	// Get command/first word from commandline
	var err error
	err = recutil.Validate(r)
	if err != nil {
		rec.Errors = append(rec.Errors, "Validate error:"+err.Error())
		// rec, _ := record.ToString()
		// sugar.Println("Invalid command:", rec)
		rec.Invalid = true
	}
	rec.Command, rec.FirstWord, err = GetCommandAndFirstWord(r.CmdLine)
	if err != nil {
		rec.Errors = append(rec.Errors, "GetCommandAndFirstWord error:"+err.Error())
		// rec, _ := record.ToString()
		// sugar.Println("Invalid command:", rec)
		rec.Invalid = true // should this be really invalid ?
	}
	return rec
}
