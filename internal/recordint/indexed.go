package recordint

import "github.com/curusarn/resh/record"

// Indexed record allows us to find records in history file in order to edit them
type Indexed struct {
	Rec record.V1
	Idx int
}
