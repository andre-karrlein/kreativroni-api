package model

type Customer struct {
	ID             string `json:"id" firestore:"id,omitempty"`
	CustomerNumber string `json:"customerNumber" firestore:"customerNumber,omitempty"`
	Name           string `json:"name" firestore:"name,omitempty"`
	Surname        string `json:"surname" firestore:"surname,omitempty"`
	UserOrderId    string `json:"userOrderId" firestore:"userOrderId,omitempty"`
}
