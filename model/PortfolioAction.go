package model

import "time"

type PortfolioAction struct {
	Name     string
	Quantity float64
	Date     time.Time
	Price    float64
}
