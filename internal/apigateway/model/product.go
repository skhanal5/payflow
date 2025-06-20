package model

type ProductResponse struct {
	ID 		string  `json:"id"`	
	Name 	string  `json:"name"`
	Description string `json:"description"`
	Price 	float64 `json:"price"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
	Stock 	int     `json:"stock"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int                `json:"total"`
}