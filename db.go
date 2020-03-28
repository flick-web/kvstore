package kvstore

import "errors"

// The KeyValueStore interface represents a type that supports basic set, get,
// and delete methods on arbitrary values. Implementers can back the store with
// any method (local file, external database, in-memory, etc). Therefore, there
// is no requirement that a type that implements this interface must support
// persistence.
type KeyValueStore interface {
	SetObject(hashKey, rangeKey string, value interface{}) error
	GetObject(hashKey, rangeKey string, result interface{}) error
	DeleteObject(hashKey, rangeKey string) error
}

// ErrKeyNotFound is returned when a value could not be found for the specified key.
var ErrKeyNotFound = errors.New("Value not found for key")
