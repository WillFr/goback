package lib

import (
	"fmt"
	"sort"

	"github.com/willfr/goback/model"
)

func Verify(history []model.PortfolioAction, inputs *model.StrategyInputs) float64 {
	sort.Slice(history, func(i, j int) bool {
		return history[i].Date.Before(history[j].Date)
	})
	gain := 0.0
	capital := (*inputs).InitialCapital
	gainPerTicker := make(map[string]float64)
	for _, action := range history {
		price:=0.0
		if action.Quantity<0 {
			price = action.Low
		} else {
			price = action.Price
		}
		diff := -action.Quantity * price
		gain += diff
		capital += diff
		if capital < 0{
			fmt.Println(action.Date, "  ", action, ">>>", diff, " : ", capital )
		}
		gainPerTicker[action.Name] += diff
	}
	fmt.Println("verified gain: ", gain)
	type Pair struct {
		key string
		val float64
	}
	tickerList := make([]Pair, len(gainPerTicker))
	i := 0
	for key, value := range gainPerTicker {
		tickerList[i] = Pair{key: key, val: value}
		i++
	}
	sort.Slice(tickerList, func(i, j int) bool {
		return tickerList[i].val < tickerList[j].val
	})
	fmt.Println("Bottom 10: ")
	for i := 0; i < min(len(tickerList), 10); i++ {
		fmt.Println(i, " ", tickerList[i].key, ": ", tickerList[i].val)
	}

	fmt.Println("Top 10: ")
	L := len(tickerList) - 1
	for i := 0; i < min(len(tickerList), 10); i++ {
		fmt.Println(i, " ", tickerList[L-i].key, ": ", tickerList[L-i].val)
	}
	return gain
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
