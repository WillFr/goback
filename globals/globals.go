package globals

import (
	"sync"

	"github.com/willfr/goback/model"
)

var Portfolio = make(map[string]float64)
var History = make([]model.PortfolioAction, 0)
var Capital float64
var Invested float64
var Total float64
var OpCount int
var Mutex = &sync.Mutex{}

var Stoped int
var Gained int
var MarketClosed int
