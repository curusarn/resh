package records

// DEPRECATION NOTICE: This package should be removed in favor of:
// - record: public record definitions
// - recordint: internal record definitions
// - recutil: record-related utils

import (
	"bufio"
	"os"
	"strings"

	"github.com/curusarn/resh/internal/histlist"
	"go.uber.org/zap"
)

// LoadCmdLinesFromZshFile loads cmdlines from zsh history file
func LoadCmdLinesFromZshFile(sugar *zap.SugaredLogger, fname string) histlist.Histlist {
	hl := histlist.New(sugar)
	file, err := os.Open(fname)
	if err != nil {
		sugar.Error("Failed to open zsh history file - skipping reading zsh history", zap.Error(err))
		return hl
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// trim newline
		line = strings.TrimRight(line, "\n")
		var cmd string
		// zsh format EXTENDED_HISTORY
		// : 1576270617:0;make install
		// zsh format no EXTENDED_HISTORY
		// make install
		if len(line) == 0 {
			// skip empty
			continue
		}
		if strings.Contains(line, ":") && strings.Contains(line, ";") &&
			len(strings.Split(line, ":")) >= 3 && len(strings.Split(line, ";")) >= 2 {
			// contains at least 2x ':' and 1x ';' => assume EXTENDED_HISTORY
			cmd = strings.Split(line, ";")[1]
		} else {
			cmd = line
		}
		hl.AddCmdLine(cmd)
	}
	return hl
}

// LoadCmdLinesFromBashFile loads cmdlines from bash history file
func LoadCmdLinesFromBashFile(sugar *zap.SugaredLogger, fname string) histlist.Histlist {
	hl := histlist.New(sugar)
	file, err := os.Open(fname)
	if err != nil {
		sugar.Error("Failed to open bash history file - skipping reading bash history", zap.Error(err))
		return hl
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// trim newline
		line = strings.TrimRight(line, "\n")
		// trim spaces from left
		line = strings.TrimLeft(line, " ")
		// bash format (two lines)
		// #1576199174
		// make install
		if strings.HasPrefix(line, "#") {
			// is either timestamp or comment => skip
			continue
		}
		if len(line) == 0 {
			// skip empty
			continue
		}
		hl.AddCmdLine(line)
	}
	return hl
}
