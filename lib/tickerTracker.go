package lib

import (
	"os"
	"bufio"
	"container/list"
	"math"
	"time"
	//"fmt"

	"github.com/willfr/goback/globals"
	"github.com/willfr/goback/model"
	"github.com/willfr/goback/model/action"
	"github.com/willfr/goback/model/reason"
)

func TrackTicker(ticker string, strategy *model.StrategyInputs, input chan time.Time, output chan time.Time) {
	windowQ := list.New()
	volumeSum := 0.0
	var stashed *model.DataPoint = nil
	bought := float64(0.0)
	boughtAt := float64(0.0)
	lastBought := 0
	var currentTime time.Time

	binFilePath := GetTickerFilePath(ticker) + ".bin"
	_, err := os.Stat(binFilePath)
	parseFromBin := false
	var writer *bufio.Writer
	var reader *bufio.Reader
	var scanner *bufio.Scanner
	if os.IsNotExist(err) {
		scanner = bufio.NewScanner(OpenTicker(ticker))
		f,_ := os.OpenFile(binFilePath, os.O_CREATE, 0600)
		writer = bufio.NewWriter(f)
		defer writer.Flush()
	} else {
		parseFromBin = true
		f,_ := os.Open(binFilePath)
		reader = bufio.NewReader(f)
	}

	lineI := 0
	for {
		lineI++
		if stashed == nil {
			if !parseFromBin {
				if !scanner.Scan() {
					break
				}
				line := scanner.Bytes()
				dp := ParseLine(line)
				if dp == nil {
					continue
				}
				//data, err := bson.Marshal(&d)
				
				WriteToBin(writer, dp)
				if lineI % 1000 == 0 {
					writer.Flush()
				}
				stashed = dp
			} else{
				var err error
				stashed, err =ReadFromBin(reader)
				if err != nil {
					break
				}
			}
		}
		currentTime = <-input
		if (*stashed).Date == currentTime {
			current := *stashed
			stashed = nil

			windowQ.PushBack(current)
			volumeSum += current.Volume

			if windowQ.Len() > (*strategy).Window {

				toRemove := windowQ.Remove(windowQ.Front()).(model.DataPoint)
				volumeSum -= toRemove.Volume
				revenue := 0.0
				gain := 0.0

				_action, _reason := (*strategy).Decide(windowQ, &current, bought, boughtAt, lastBought, volumeSum)
				if _action != action.NONE {
					if _action == action.SOLD {
						gain = bought*(current.Low-boughtAt) - (*strategy).Commission
						//tax := 0.0
						//if gain > 0 {
						//	tax = gain * 0.00
						//}
						revenue = bought*current.Low /*- tax*/ - (*strategy).Commission

						globals.Mutex.Lock()
						globals.Invested -= bought * boughtAt
						globals.Capital += revenue
						globals.Total += gain

						switch _reason {
						case reason.MARKET_CLOSED:
							globals.MarketClosed++
						case reason.GAINED:
							globals.Gained++
						case reason.STOP:
							globals.Stoped++
						}
						delete(globals.Portfolio, ticker)
						globals.History = append(globals.History, model.PortfolioAction{Date: current.Date, Quantity: -bought, Name: ticker, Price: current.Price, Low: current.Low})
						globals.Mutex.Unlock()

						bought = 0
					} else if _action == action.BOUGHT {
						globals.Mutex.Lock()
						batch := strategy.Batch
						if globals.Capital/8 > batch {
							batch = globals.Capital / 8
						}
						if globals.Capital > batch {
							prevBought := bought
							bought = math.Floor(batch / current.Price)
							globals.Portfolio[ticker] += bought
							paid := bought * current.Price
							globals.Invested += paid
							boughtAt = (current.Price*bought + prevBought*boughtAt) / (bought + prevBought)
							revenue = -paid - strategy.Commission
							lastBought = current.Date.Day()
							globals.Capital += revenue
							globals.OpCount++
							globals.History = append(globals.History, model.PortfolioAction{Date: current.Date, Quantity: bought, Name: ticker, Price: current.Price})
							bought += prevBought

						}
						globals.Mutex.Unlock()
					}
				}
			}
			output <- currentTime.Add(time.Minute * time.Duration(1))
		} else {
			output <- (*stashed).Date
		}
	}

	if bought > 0 {

		globals.Mutex.Lock()
		globals.Invested -= bought * boughtAt
		globals.Capital += bought * boughtAt

		delete(globals.Portfolio, ticker)
		globals.History = append(globals.History, model.PortfolioAction{Date: currentTime, Quantity: -bought, Name: ticker, Price: boughtAt})

		globals.Mutex.Unlock()
	}
	// receive one last time because o the loop below
	<-input
	close(output)
}
