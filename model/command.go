package model

type PaymentDebitedCommand struct {
	OrderID    string
	CustomerID string
	Amount     float64
}
