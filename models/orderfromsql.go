package models

type OrderFromDB struct {
	OrderID        int
	OrderTime      string
	OrderElemCount int
	OrderValue     int
	RestID         int
	UserID         int
	ItemID         []uint8
}
