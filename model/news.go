package model

type News struct {
	ID        string `json:"id"`
	Sort_key  string `json:"sort_key"`
	Title     string `json:"title"`
	Teaser    string `json:"teaser"`
	Message   string `json:"message"`
	Available bool   `json:"available"`
}
