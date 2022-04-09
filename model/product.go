package model

type Listing struct {
	Id int `json:"listing_id"`
}

type ListingProduct struct {
	Id            int            `json:"listing_id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Quantity      int            `json:"quantity"`
	ShopSectionId int            `json:"shop_section_id"`
	Tags          []string       `json:"tags"`
	Price         Price          `json:"price"`
	Url           string         `json:"url"`
	Images        []ProductImage `json:"images"`
}
type Product struct {
	Id            int            `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Quantity      int            `json:"quantity"`
	ShopSectionId int            `json:"shop_section_id"`
	Tags          []string       `json:"tags"`
	Price         Price          `json:"price"`
	Url           string         `json:"url"`
	Images        []ProductImage `json:"images"`
}

type ProductImage struct {
	Url_full    string `json:"url_fullxfull"`
	Url_75x75   string `json:"url_75x75"`
	Url_170x135 string `json:"url_170x135"`
	Url_570xN   string `json:"url_570xN"`
}

type Price struct {
	Amount        int    `json:"amount"`
	Divisor       int    `json:"divisor"`
	Currency_Code string `json:"currency_code"`
}

type EtsyProductData struct {
	Count   int              `json:"count"`
	Results []ListingProduct `json:"results"`
}

type ListingData struct {
	Count   int       `json:"count"`
	Results []Listing `json:"results"`
}
