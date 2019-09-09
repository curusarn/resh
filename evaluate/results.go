package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"strconv"
)

type results struct {
	writer                  *bufio.Writer
	size                    int
	matches                 []int // matches[N] -> # of matches at distance N
	matchesTotal            int
	charactersRecalled      []int
	charactersRecalledTotal int
	dataPointCount          int
}

func (r *results) init() {
	r.matches = make([]int, r.size)
	r.charactersRecalled = make([]int, r.size)
}

func (r *results) addMatch(distance int, cmdLength int) {
	if distance >= r.size {
		// --calculate-total
		// log.Fatal("Match distance is greater than size of statistics")
		r.matchesTotal++
		r.charactersRecalledTotal += cmdLength
		return
	}
	r.matches[distance]++
	r.matchesTotal++
	r.charactersRecalled[distance] += cmdLength
	r.charactersRecalledTotal += cmdLength
	r.dataPointCount++
}

func (r *results) addMiss() {
	r.dataPointCount++
}

func (r *results) printCumulative() {
	matchesPercent := 0.0
	out := "### Matches ###\n"
	for i := 0; i < r.size; i++ {
		matchesPercent += 100 * float64(r.matches[i]) / float64(r.dataPointCount)
		out += strconv.Itoa(i) + " ->"
		out += fmt.Sprintf(" (%.1f %%)\n", matchesPercent)
		for j := 0; j < int(math.Round(matchesPercent)); j++ {
			out += "#"
		}
		out += "\n"
	}
	matchesPercent = 100 * float64(r.matchesTotal) / float64(r.dataPointCount)
	out += "TOTAL ->"
	out += fmt.Sprintf(" (%.1f %%)\n", matchesPercent)
	for j := 0; j < int(math.Round(matchesPercent)); j++ {
		out += "#"
	}
	out += "\n"

	n, err := r.writer.WriteString(string(out) + "\n\n")
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("Nothing was written", n)
	}

	charsRecall := 0.0
	out = "### Characters recalled per submission ###\n"
	for i := 0; i < r.size; i++ {
		charsRecall += float64(r.charactersRecalled[i]) / float64(r.dataPointCount)
		out += strconv.Itoa(i) + " ->"
		out += fmt.Sprintf(" (%.2f)\n", charsRecall)
		for j := 0; j < int(math.Round(charsRecall)); j++ {
			out += "#"
		}
		out += "\n"
	}
	charsRecall = float64(r.charactersRecalledTotal) / float64(r.dataPointCount)
	out += "TOTAL ->"
	out += fmt.Sprintf(" (%.2f)\n", charsRecall)
	for j := 0; j < int(math.Round(charsRecall)); j++ {
		out += "#"
	}
	out += "\n"

	n, err = r.writer.WriteString(string(out) + "\n\n")
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("Nothing was written", n)
	}
}
