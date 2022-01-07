package models

type OrderFromClient struct {
	OrderID   int
	SentBy    string
	OrderTime string
	ItemVal   float64
	ClosedBy  string
	RestName  string
	PaidType  string
	From      string
	Item      string
	ItemCount int
	TableName string
	ClosedAt  string
}
