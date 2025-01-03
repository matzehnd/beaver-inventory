package domain

type Quantity struct {
	Amount float64 `json:"amount"`
	Unit   Unit    `json:"unit"`
}

type Unit string

type StockItemKey struct {
	LocationId string
	BatchId    string
}

type StockItem struct {
	LocationId string   `json:"locationId"`
	BatchId    string   `json:"batchId"`
	Quantity   Quantity `json:"quantity"`
}
