package domain

type Product struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProductStock struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Quantity Quantity `json:"quantity"`
}
