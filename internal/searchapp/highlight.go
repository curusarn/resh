package searchapp

import (
	"strconv"
	"strings"
)

func cleanHighlight(str string) string {
	prefix := "\033["

	invert := "\033[7;1m"
	invertGreen := "\033[32;7;1m"
	end := "\033[0m"
	replace := []string{invert, invertGreen, end}
	for i := 30; i < 48; i++ {
		base := prefix + strconv.Itoa(i)
		normal := base + "m"
		bold := base + ";1m"
		replace = append(replace, normal, bold)
	}
	if strings.Contains(str, prefix) == false {
		return str
	}
	for _, escSeq := range replace {
		str = strings.ReplaceAll(str, escSeq, "")
	}
	return str
}

func highlightHeader(str string) string {
	underline := "\033[4m"
	end := "\033[0m"
	// no clean highlight
	return underline + str + end
}

func highlightStatus(str string) string {
	invert := "\033[7;1m"
	end := "\033[0m"
	return invert + cleanHighlight(str) + end
}

func highlightSelected(str string) string {
	// template "\033[3%d;%dm"
	// invertGreen := "\033[32;7;1m"
	invert := "\033[7;1m"
	end := "\033[0m"
	return invert + cleanHighlight(str) + end
}

func highlightDate(str string) string {
	// template "\033[3%d;%dm"
	yellowNormal := "\033[33m"
	end := "\033[0m"
	return yellowNormal + cleanHighlight(str) + end
}

func highlightHost(str string) string {
	// template "\033[3%d;%dm"
	redNormal := "\033[31m"
	end := "\033[0m"
	return redNormal + cleanHighlight(str) + end
}

func highlightPwd(str string) string {
	// template "\033[3%d;%dm"
	blueBold := "\033[34;1m"
	end := "\033[0m"
	return blueBold + cleanHighlight(str) + end
}

func highlightMatch(str string) string {
	// template "\033[3%d;%dm"
	magentaBold := "\033[35;1m"
	end := "\033[0m"
	return magentaBold + cleanHighlight(str) + end
}

func highlightWarn(str string) string {
	// template "\033[3%d;%dm"
	// orangeBold := "\033[33;1m"
	redBold := "\033[31;1m"
	end := "\033[0m"
	return redBold + cleanHighlight(str) + end
}

func highlightGit(str string) string {
	// template "\033[3%d;%dm"
	greenBold := "\033[32;1m"
	end := "\033[0m"
	return greenBold + cleanHighlight(str) + end
}

// DoHighlightHeader .
func DoHighlightHeader(str string, minLength int) string {
	if len(str) < minLength {
		str = str + strings.Repeat(" ", minLength-len(str))
	}
	return highlightHeader(str)
}

// DoHighlightString .
func DoHighlightString(str string, minLength int) string {
	if len(str) < minLength {
		str = str + strings.Repeat(" ", minLength-len(str))
	}
	return highlightSelected(str)
}

// EXTRAS

func highlightModeTitle(str string) string {
	// template "\033[3%d;%dm"
	greenNormal := "\033[32;1m"
	end := "\033[0m"
	return greenNormal + cleanHighlight(str) + end
}
