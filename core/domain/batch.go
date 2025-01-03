package domain

import "time"

type Batch struct {
	Id           string    `json:"id"`
	Product      Product   `json:"product"`
	ProducedAt   time.Time `json:"producedAt"`
	SellLatestAt time.Time `json:"sellLatestAt"`
}
