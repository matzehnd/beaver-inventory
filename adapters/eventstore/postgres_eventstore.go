package eventstore

import (
	"beaver/inventory/core/domain"
	"database/sql"
	"encoding/json"
	"fmt"
)

type PostgresEventStore struct {
	db *sql.DB
}

func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

func (es *PostgresEventStore) Save(event interface{}) error {
	var eventType string
	switch e := event.(type) {
	case domain.StockChangeEvent:
		eventType = "stock-change"

	default:
		return fmt.Errorf("unknown event type: %T", e)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = es.db.Exec("INSERT INTO inventory_events (type, data) VALUES ($1,$2)", eventType, data)
	return err
}

func (es *PostgresEventStore) Load() ([]interface{}, error) {
	rows, err := es.db.Query("SELECT type, data FROM inventory_events order by id ASC")
	if err != nil {
		return nil, fmt.Errorf("unable to select events: %w", err)
	}
	defer rows.Close()

	var events []interface{}
	for rows.Next() {
		var eventType string
		var data []byte
		if err := rows.Scan(&eventType, &data); err != nil {
			return nil, err
		}
		switch eventType {
		case "stock-change":
			var event domain.StockChangeEvent
			if err := json.Unmarshal(data, &event); err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}
