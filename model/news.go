package model

type News struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Teaser    string `json:"teaser"`
	Message   string `json:"message"`
	Available bool   `json:"available"`
}
