package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"strconv"
)

type statistics struct {
	writer                  *bufio.Writer
	size                    int
	matches                 []int
	matchesTotal            int
	charactersRecalled      []int
	charactersRecalledTotal int
	dataPointCount          int
}

func (s *statistics) init() {
	s.matches = make([]int, s.size)
	s.charactersRecalled = make([]int, s.size)
}

func (s *statistics) addMatch(distance int, cmdLength int) {
	if distance >= s.size {
		// --calculate-total
		// log.Fatal("Match distance is greater than size of statistics")
		s.matchesTotal++
		s.charactersRecalledTotal += cmdLength
		return
	}
	s.matches[distance]++
	s.matchesTotal++
	s.charactersRecalled[distance] += cmdLength
	s.charactersRecalledTotal += cmdLength
	s.dataPointCount++
}

func (s *statistics) addMiss() {
	s.dataPointCount++
}

func (s *statistics) printCumulative() {
	matchesPercent := 0.0
	out := "### Matches ###\n"
	for i := 0; i < s.size; i++ {
		matchesPercent += 100 * float64(s.matches[i]) / float64(s.dataPointCount)
		out += strconv.Itoa(i) + " ->"
		out += fmt.Sprintf(" (%.1f %%)\n", matchesPercent)
		for j := 0; j < int(math.Round(matchesPercent)); j++ {
			out += "#"
		}
		out += "\n"
	}
	matchesPercent = 100 * float64(s.matchesTotal) / float64(s.dataPointCount)
	out += "TOTAL ->"
	out += fmt.Sprintf(" (%.1f %%)\n", matchesPercent)
	for j := 0; j < int(math.Round(matchesPercent)); j++ {
		out += "#"
	}
	out += "\n"

	n, err := s.writer.WriteString(string(out) + "\n\n")
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("Nothing was written", n)
	}

	charsRecall := 0.0
	out = "### Characters recalled per submission ###\n"
	for i := 0; i < s.size; i++ {
		charsRecall += float64(s.charactersRecalled[i]) / float64(s.dataPointCount)
		out += strconv.Itoa(i) + " ->"
		out += fmt.Sprintf(" (%.2f)\n", charsRecall)
		for j := 0; j < int(math.Round(charsRecall)); j++ {
			out += "#"
		}
		out += "\n"
	}
	charsRecall = float64(s.charactersRecalledTotal) / float64(s.dataPointCount)
	out += "TOTAL ->"
	out += fmt.Sprintf(" (%.2f)\n", charsRecall)
	for j := 0; j < int(math.Round(charsRecall)); j++ {
		out += "#"
	}
	out += "\n"

	n, err = s.writer.WriteString(string(out) + "\n\n")
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("Nothing was written", n)
	}
}
