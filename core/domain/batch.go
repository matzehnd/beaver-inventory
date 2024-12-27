package domain

import "time"

type Batch struct {
	Id           string
	Product      Product
	ProducedAt   time.Time
	SellLatestAt time.Time
}
