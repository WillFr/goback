package model

import "time"

type DataPoint struct {
	Date   time.Time
	Price  float64
	Low    float64
	Volume float64
}
