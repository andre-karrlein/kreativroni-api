package model

type News struct {
	ID        string `json:"id" firestore:"id,omitempty"`
	Title     string `json:"title" firestore:"title,omitempty"`
	Teaser    string `json:"teaser" firestore:"teaser,omitempty"`
	Message   string `json:"message" firestore:"message,omitempty"`
	Available bool   `json:"available" firestore:"available,omitempty"`
}
