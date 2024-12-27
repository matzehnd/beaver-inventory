package domain

type StockChangeEvent struct {
	Batch    Batch
	Location Location
	Quantity Quantity
}
