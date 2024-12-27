package domain

type EventStore interface {
	Save(event interface{}) error
	Load() ([]interface{}, error)
}
