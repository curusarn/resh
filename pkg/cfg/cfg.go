package cfg

// Config struct
type Config struct {
	Port                    int
	SesswatchPeriodSeconds  uint
	SesshistInitHistorySize int
	Debug                   bool
	BindArrowKeysBash       bool
	BindArrowKeysZsh        bool
}
