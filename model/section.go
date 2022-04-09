package model

type SectionListing struct {
	Id    int    `json:"shop_section_id"`
	Title string `json:"title"`
}

type Section struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type SectionData struct {
	Count   int              `json:"count"`
	Results []SectionListing `json:"results"`
}
