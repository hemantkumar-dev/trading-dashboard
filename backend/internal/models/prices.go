package models

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
