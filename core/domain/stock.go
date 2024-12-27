package domain

type Quantity struct {
	Amount float64
	Unit   Unit
}

type Unit string

type StockItemKey struct {
	LocationId string
	BatchId    string
}

type StockItem struct {
	LocationId string
	BatchId    string
	Quantity   Quantity
}
