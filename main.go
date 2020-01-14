package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"gonum.org/v1/plot/plotter"

	"github.com/willfr/goback/globals"
	"github.com/willfr/goback/lib"
	"github.com/willfr/goback/model"
)

func main() {
	capitalP := flag.Float64("cap", 100000, "The initial capital")              //float64(100000)
	thresholdP := flag.Float64("thre", 1.03, "The buy threshold")               //float64(1.03)
	volumeThresholdP := flag.Float64("vthre", 1.03, "The buy volume threshold") //float64(1.03)
	windowP := flag.Int("win", 60, "The window length in minutes")              //60 //minutes
	batchP := flag.Float64("bat", 10000, "The batch size in $")
	limitP := flag.Float64("lim", 1.003, "The limit ratio")
	stopP := flag.Float64("sto", 0.97, "The stop ratio")
	commissionP := flag.Float64("com", 1, "The commission per transaction")
	cutoffP := flag.Int("cut", 13, "The cuttoff hour")
	startHP := flag.Int("sta", 9, "The start hour")
	startHMP := flag.Int("stam", 30, "The start minute")
	maxFileP := flag.Int("maxFile", 3, "The number of file to be parsed")
	minVolumeSumP := flag.Float64("mvs", 0.0, "The minimum volume sum in the window to buy")
	randSeedP := flag.Int64("rs", time.Now().UnixNano(), "The seed for the randomizer")
	flag.Parse()

	strategyInputs := model.StrategyInputs{
		InitialCapital:  *capitalP,
		Threshold:       *thresholdP,
		VolumeThreshold: *volumeThresholdP,
		Window:          *windowP,
		Batch:           *batchP,
		Limit:           *limitP,
		Stop:            *stopP,
		Commission:      *commissionP,
		Cutoff:          uint8(*cutoffP),
		StartH:          uint8(*startHP),
		StartHM:         uint8(*startHMP),
		MinVolumeSum:    *minVolumeSumP}

	maxFile := *maxFileP
	randSeed := *randSeedP
	globals.Total = float64(0)
	globals.OpCount = 0
	globals.Capital = strategyInputs.InitialCapital
	MAX_TIME := model.SimplifiedDate{Year: 2070}

	tickers := lib.ListTickers("C:\\Users\\Guillaume\\Desktop\\stocks\\intraday\\")
	//shuffle the files
	rand.Seed(randSeed)
	rand.Shuffle(len(tickers), func(i, j int) { tickers[i], tickers[j] = tickers[j], tickers[i] })

	clock := model.SimplifiedDate{Year: 1997, Month: 01, Day: 01}

	nbChan := 0
	inputs := make([]chan byte, 600)
	closedChan := make([]bool, 600)

	gainPts := plotter.XYs{}
	capitalPts := plotter.XYs{}
	opPts := plotter.XYs{}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		lib.DrawGraph(gainPts, capitalPts, opPts)
		os.Exit(1)
	}()

	for _, ticker := range tickers {
		fmt.Println(ticker)
		inputs[nbChan] = make(chan model.SimplifiedDate)
		go lib.TrackTicker(ticker, &strategyInputs, inputs[nbChan])
		inputs[nbChan] <- clock

		nbChan++
		if nbChan == maxFile {
			break
		}
	}

	nbClosed := 0
	dpI := 0
	for nbClosed < nbChan {
		min := MAX_TIME

		for i := 0; i < nbChan; i++ {
			if closedChan[i] {
				continue
			}
			t, open := <-inputs[i]

			if open {
				if t.Before(min) {
					min = t
				} else {

				}
			} else if !open {
				closedChan[i] = true
				nbClosed++
			}
		}

		clock = min
		for i := 0; i < nbChan; i++ {
			if closedChan[i] {
				continue
			}
			inputs[i] <- clock
		}
		if dpI%120000 == 0 {
			fmt.Println(clock)
			fmt.Printf("GAIN: %.2f \n", globals.Total)
			fmt.Printf("CAPITAL: %.2f \n", globals.Capital)
			fmt.Printf("INVESTED: %.2f \n", globals.Invested)
			fmt.Println("OP: ", globals.OpCount, " STOP: ", globals.Stoped, " GAINED: ", globals.Gained, " CLOSED: ", globals.MarketClosed)
			fmt.Println(globals.Portfolio)
			fmt.Println()

		}
		if dpI%5000 == 0 {
			gainPts = append(gainPts, plotter.XY{float64(clock.Unix()), globals.Total})
			capitalPts = append(capitalPts, plotter.XY{float64(clock.Unix()), globals.Capital})
			opPts = append(opPts, plotter.XY{float64(clock.Unix()), float64(globals.OpCount)})
		}

		dpI++
	}

	fmt.Printf("GAIN: %.2f \n", globals.Total)
	fmt.Printf("CAPITAL: %.2f \n", globals.Capital)
	fmt.Printf("INVESTED: %.2f \n", globals.Invested)
	lib.DrawGraph(gainPts, capitalPts, opPts)
	lib.SaveRun()
	lib.GenerateHistory()

	fmt.Println("Verifying :")
	fmt.Println(lib.Verify(globals.History, &strategyInputs))

	//os.Exit(int(globals.Total))
}
