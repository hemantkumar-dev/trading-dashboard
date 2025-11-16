package models

type Order struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
