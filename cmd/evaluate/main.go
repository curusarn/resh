package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/curusarn/resh/pkg/histanal"
	"github.com/curusarn/resh/pkg/records"
	"github.com/curusarn/resh/pkg/strat"
)

// version from git set during build
var version string

// commit from git set during build
var commit string

func main() {
	const maxCandidates = 50

	usr, _ := user.Current()
	dir := usr.HomeDir
	historyPath := filepath.Join(dir, ".resh_history.json")
	historyPathBatchMode := filepath.Join(dir, "resh_history.json")
	sanitizedHistoryPath := filepath.Join(dir, "resh_history_sanitized.json")
	// tmpPath := "/tmp/resh-evaluate-tmp.json"

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")
	input := flag.String("input", "",
		"Input file (default: "+historyPath+" OR "+sanitizedHistoryPath+
			" depending on --sanitized-input option)")
	// outputDir := flag.String("output", "/tmp/resh-evaluate", "Output directory")
	sanitizedInput := flag.Bool("sanitized-input", false,
		"Handle input as sanitized (also changes default value for input argument)")
	plottingScript := flag.String("plotting-script", "resh-evaluate-plot.py", "Script to use for plotting")
	inputDataRoot := flag.String("input-data-root", "",
		"Input data root, enables batch mode, looks for files matching --input option")
	slow := flag.Bool("slow", false,
		"Enables strategies that takes a long time (e.g. markov chain strategies).")
	skipFailedCmds := flag.Bool("skip-failed-cmds", false,
		"Skips records with non-zero exit status.")
	debugRecords := flag.Float64("debug", 0, "Debug records - percentage of records that should be debugged.")

	flag.Parse()

	// handle show{Version,Revision} options
	if *showVersion == true {
		fmt.Println(version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(commit)
		os.Exit(0)
	}

	// handle batch mode
	batchMode := false
	if *inputDataRoot != "" {
		batchMode = true
	}
	// set default input
	if *input == "" {
		if *sanitizedInput {
			*input = sanitizedHistoryPath
		} else if batchMode {
			*input = historyPathBatchMode
		} else {
			*input = historyPath
		}
	}

	var evaluator histanal.HistEval
	if batchMode {
		evaluator = histanal.NewHistEvalBatchMode(*input, *inputDataRoot, maxCandidates, *skipFailedCmds, *debugRecords, *sanitizedInput)
	} else {
		evaluator = histanal.NewHistEval(*input, maxCandidates, *skipFailedCmds, *debugRecords, *sanitizedInput)
	}

	var simpleStrategies []strat.ISimpleStrategy
	var strategies []strat.IStrategy

	// dummy := strategyDummy{}
	// simpleStrategies = append(simpleStrategies, &dummy)

	simpleStrategies = append(simpleStrategies, &strat.Recent{})

	// frequent := strategyFrequent{}
	// frequent.init()
	// simpleStrategies = append(simpleStrategies, &frequent)

	// random := strategyRandom{candidatesSize: maxCandidates}
	// random.init()
	// simpleStrategies = append(simpleStrategies, &random)

	directory := strat.DirectorySensitive{}
	directory.Init()
	simpleStrategies = append(simpleStrategies, &directory)

	// dynamicDistG := strat.DynamicRecordDistance{
	// 	MaxDepth:   3000,
	// 	DistParams: records.DistParams{Pwd: 10, RealPwd: 10, SessionID: 1, Time: 1, Git: 10},
	// 	Label:      "10*pwd,10*realpwd,session,time,10*git",
	// }
	// dynamicDistG.Init()
	// strategies = append(strategies, &dynamicDistG)

	distanceStaticBest := strat.RecordDistance{
		MaxDepth:   3000,
		DistParams: records.DistParams{Pwd: 10, RealPwd: 10, SessionID: 1, Time: 1},
		Label:      "10*pwd,10*realpwd,session,time",
	}
	strategies = append(strategies, &distanceStaticBest)

	recentBash := strat.RecentBash{}
	recentBash.Init()
	strategies = append(strategies, &recentBash)

	if *slow {

		markovCmd := strat.MarkovChainCmd{Order: 1}
		markovCmd.Init()

		markovCmd2 := strat.MarkovChainCmd{Order: 2}
		markovCmd2.Init()

		markov := strat.MarkovChain{Order: 1}
		markov.Init()

		markov2 := strat.MarkovChain{Order: 2}
		markov2.Init()

		simpleStrategies = append(simpleStrategies, &markovCmd2, &markovCmd, &markov2, &markov)
	}

	for _, strategy := range simpleStrategies {
		strategies = append(strategies, strat.NewSimpleStrategyWrapper(strategy))
	}

	for _, strat := range strategies {
		err := evaluator.Evaluate(strat)
		if err != nil {
			log.Println("Evaluator evaluate() error:", err)
		}
	}

	evaluator.CalculateStatsAndPlot(*plottingScript)
}
