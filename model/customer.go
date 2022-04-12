package model

type Customer struct {
	ID             string `json:"id"`
	Sort_key       string `json:"sort_key"`
	CustomerNumber string `json:"customerNumber"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	UserOrderId    string `json:"userOrderId"`
}
