package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"sort"

	"github.com/wcharczuk/go-chart"
)

type statistics struct {
	//size                    int
	dataPointCount int
	cmdLineCount   map[string]int
}

func (s *statistics) init() {
	s.cmdLineCount = make(map[string]int)
}

func (s *statistics) addCmdLine(cmdLine string, cmdLength int) {
	s.cmdLineCount[cmdLine]++
	s.dataPointCount++
}

func (s *statistics) graphCmdFrequencyAsFuncOfRank() {

	var xValues []float64
	var yValues []float64

	sortedValues := sortMapByvalue(s.cmdLineCount)
	sortedValues = sortedValues[:100] // cut off at rank 100

	normalizeCoeficient := float64(s.dataPointCount) / float64(sortedValues[0].Value)
	for i, pair := range sortedValues {
		rank := i + 1
		frequency := float64(pair.Value) / float64(s.dataPointCount)
		normalizeFrequency := frequency * normalizeCoeficient

		xValues = append(xValues, float64(rank))
		yValues = append(yValues, normalizeFrequency)
	}

	graphName := "cmdFrqAsFuncOfRank"
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(), //enables / displays the x-axis
			Ticks: []chart.Tick{
				{0.0, "0"},
				{1.0, "1"},
				{2.0, "2"},
				{3.0, "3"},
				{4.0, "4"},
				{5.0, "5"},
				{10.0, "10"},
				{15.0, "15"},
				{20.0, "20"},
				{25.0, "25"},
				{30.0, "30"},
				{35.0, "35"},
				{40.0, "40"},
				{45.0, "45"},
				{50.0, "50"},
			},
		},
		YAxis: chart.YAxis{
			AxisType: chart.YAxisSecondary,
			Style:    chart.StyleShow(), //enables / displays the y-axis
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
					DotColor:    chart.GetDefaultColor(0),
					DotWidth:    3.0,
				},
				XValues: xValues,
				YValues: yValues,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Fatal("chart.Render error:", err)
	}
	ioutil.WriteFile("/tmp/resh-graph_"+graphName+".png", buffer.Bytes(), 0644)
}

func sortMapByvalue(input map[string]int) []Pair {
	p := make(PairList, len(input))

	i := 0
	for k, v := range input {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}

// Pair - A data structure to hold key/value pairs
type Pair struct {
	Key   string
	Value int
}

// PairList - A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
