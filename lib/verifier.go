package lib

import (
	"fmt"
	"sort"

	"github.com/willfr/goback/model"
)

func Verify(history []model.PortfolioAction) float64 {
	sort.Slice(history, func(i, j int) bool {
		return history[i].Name < history[j].Name && history[i].Date.Before(history[j].Date)
	})
	gain := 0.0
	gainPerTicker := make(map[string]float64)
	for _, action := range history {
		diff := -action.Quantity * action.Price
		gain += diff
		gainPerTicker[action.Name] += diff
	}

	type Pair struct {
		a, b interface{}
	}
	tickerList := make([]Pair, len(gainPerTicker))
	i := 0
	for key, value := range gainPerTicker {
		tickerList[i] = Pair{a: key, b: value}
	}
	sort.Slice(tickerList, func(i, j int) bool {
		return tickerList[i].b.(float64) < tickerList[j].b.(float64)
	})
	fmt.Println("Bottom 10: ")
	for i := 0; i < min(len(tickerList), 10); i++ {
		fmt.Println(tickerList[i].a, ": ", tickerList[i].b)
	}

	fmt.Println("Top 10: ")
	L := len(tickerList) - 1
	for i := 0; i < min(len(tickerList), 10); i++ {
		fmt.Println(tickerList[L-i].a, ": ", tickerList[L-i].b)
	}
	return gain
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
