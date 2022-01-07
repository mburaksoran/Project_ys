package models

type DataForClient struct {
	OrderTime    string  `json:"time"`
	TotalPayment float64 `json:"payment"`
	DataComeFrom int     `json:"datacomefrom"`
}
