package model


type PortfolioAction struct {
	Name     string
	Quantity float64
	Date     SimplifiedDate
	Price    float64
	Low    float64
}
