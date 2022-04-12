package model

type Order struct {
	ID           string `json:"id"`
	Sort_key     string `json:"sort_key"`
	User         string `json:"user"`
	Payment      string `json:"payment"`
	ProductName  string `json:"productName"`
	ProductId    string `json:"productId"`
	Variation    string `json:"variation"`
	Price        string `json:"price"`
	Quantity     string `json:"quantity"`
	Status       string `json:"status"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AddressLine1 string `json:"addressline1"`
	AddressLine2 string `json:"addressline2"`
	City         string `json:"city"`
	InvoiceId    string `json:"invoiceId"`
}
