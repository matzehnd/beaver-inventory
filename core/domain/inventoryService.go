package domain

type StockService struct {
	eventStore EventStore
	stock      map[StockItemKey]*StockItem
	products   map[string]*Product
	locations  map[string]*Location
	batches    map[string]*Batch
}

type StockChange struct {
	BatchId    string
	LocationId string
	Quantity   float64
}

func NewStockService(eventStore EventStore) *StockService {
	return &StockService{
		eventStore: eventStore,
		stock:      make(map[StockItemKey]*StockItem),
		products:   make(map[string]*Product),
		locations:  make(map[string]*Location),
		batches:    make(map[string]*Batch),
	}
}

func (s *StockService) RebuildEventStream() error {
	events, err := s.eventStore.Load()
	if err != nil {
		return err
	}
	for _, event := range events {
		s.apply(event)
	}
	return nil
}

func (s *StockService) GetAllProducts() []*ProductStock {
	products := []*ProductStock{}
	for _, product := range s.products {
		totalStock := s.getTotalStockForProduct(product)
		products = append(products, &ProductStock{
			Id:   product.Id,
			Name: product.Name,
			Quantity: Quantity{
				Amount: totalStock,
				Unit:   "pcs",
			},
		})
	}
	return products
}

func (s *StockService) GetStock() []*StockItem {
	stockItems := []*StockItem{}
	for _, stock := range s.stock {
		stockItems = append(stockItems, stock)
	}
	return stockItems
}

func (s *StockService) GetBatches() []*Batch {
	batches := []*Batch{}
	for _, batch := range s.batches {
		batches = append(batches, batch)
	}
	return batches
}

func (s *StockService) GetLocations() []*Location {
	locations := []*Location{}
	for _, location := range s.locations {
		locations = append(locations, location)
	}
	return locations
}

func (s *StockService) StockChange(event StockChangeEvent) error {
	if err := s.eventStore.Save(event); err != nil {
		return err
	}
	s.apply(event)
	return nil
}

func (s *StockService) ApplyStockChangeEvent(event StockChangeEvent) {
	key := getInventoryItemKey(event.Location, event.Batch)
	stockItem, exists := s.stock[key]
	var newAmount float64
	if exists {
		newAmount = getZeroIfNegativ(stockItem.Quantity.Amount + event.Quantity.Amount)
	} else {
		newAmount = getZeroIfNegativ(event.Quantity.Amount)
	}
	s.stock[key] = &StockItem{
		LocationId: event.Location.Id,
		BatchId:    event.Batch.Id,
		Quantity: Quantity{
			Unit:   event.Quantity.Unit,
			Amount: newAmount,
		},
	}
	s.products[event.Batch.Product.Id] = &event.Batch.Product
	s.batches[event.Batch.Id] = &event.Batch
	s.locations[event.Location.Id] = &event.Location
}

func (s *StockService) apply(event interface{}) {
	switch e := event.(type) {
	case StockChangeEvent:
		s.ApplyStockChangeEvent(e)
	}
}
func (s *StockService) getTotalStockForProduct(product *Product) float64 {
	var sum float64
	for _, item := range s.stock {
		batch, exists := s.batches[item.BatchId]
		if exists && batch.Product.Id == product.Id {
			sum += item.Quantity.Amount
		}
	}
	return sum
}

func getInventoryItemKey(Location Location, Batch Batch) StockItemKey {
	return StockItemKey{LocationId: Location.Id, BatchId: Batch.Id}
}

func getZeroIfNegativ(amount float64) float64 {
	if amount < 0 {
		return 0
	}
	return amount
}
