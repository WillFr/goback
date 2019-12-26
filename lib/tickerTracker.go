package lib

import (
	"bufio"
	"container/list"
	"math"
	"time"

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

	scanner := bufio.NewScanner(OpenTicker(ticker))

	for {
		if stashed == nil {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			if len(line) > 20 {
				dp := ParseLine(line)
				stashed = &dp
			} else {
				continue
			}
		}
		currentTime := <-input
		if (*stashed).Date == currentTime {
			current := *stashed
			stashed = nil

			windowQ.PushBack(current)
			volumeSum += current.Volume

			if windowQ.Len() > strategy.Window {

				toRemove := windowQ.Remove(windowQ.Front()).(model.DataPoint)
				volumeSum -= toRemove.Volume
				revenue := 0.0
				gain := 0.0

				_action, _reason := strategy.Decide(windowQ, &current, bought, boughtAt, lastBought, volumeSum)
				if _action != action.NONE {
					if _action == action.SOLD {
						gain = bought*(current.Price-boughtAt) - strategy.Commission
						tax := 0.0
						if gain > 0 {
							tax = gain * 0.00
						}
						revenue = bought*current.Price - tax - strategy.Commission

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
						globals.History = append(globals.History, model.PortfolioAction{Date: current.Date, Quantity: -bought, Name: ticker, Price: current.Price})
						globals.Mutex.Unlock()

						bought = 0
					} else if _action == action.BOUGHT {
						globals.Mutex.Lock()
						if globals.Capital > strategy.Batch {
							bought = math.Floor(strategy.Batch / current.Price)
							globals.Portfolio[ticker] += bought
							paid := bought * current.Price
							globals.Invested += paid
							boughtAt = current.Price
							revenue = -paid - strategy.Commission
							lastBought = current.Date.Day()
							globals.Capital += revenue
							globals.OpCount++
							globals.History = append(globals.History, model.PortfolioAction{Date: current.Date, Quantity: bought, Name: ticker, Price: current.Price})

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
		revenue := bought*(windowQ.Back().Value.(model.DataPoint).Price) - strategy.Commission
		globals.Capital += revenue
		globals.Total += revenue - bought*boughtAt - strategy.Commission
		globals.Mutex.Unlock()
	}
	// receive one last time because of the loop below
	<-input
	close(output)
}
