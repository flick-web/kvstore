package kvstore

import (
	"bytes"
	"encoding/gob"
)

// MemoryStore is a type that fulfills the KeyValueStore interface without any
// persistence, storing values in-memory.
type MemoryStore struct {
	items map[[2]string][]byte
}

// NewMemoryStore initializes a new MemoryStore object.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[[2]string][]byte),
	}
}

// SetObject creates or updates the key-value pair.
func (kv *MemoryStore) SetObject(table, id string, value interface{}) error {
	gobBuffer := new(bytes.Buffer)
	gobEncoder := gob.NewEncoder(gobBuffer)
	err := gobEncoder.Encode(value)
	if err != nil {
		return err
	}

	key := [2]string{table, id}
	kv.items[key] = gobBuffer.Bytes()
	return nil
}

// GetObject retrieves and decodes the stored value into result.
func (kv *MemoryStore) GetObject(table, id string, result interface{}) (err error) {
	key := [2]string{table, id}
	valBytes, ok := kv.items[key]
	if !ok {
		return ErrKeyNotFound
	}

	gobBuffer := bytes.NewBuffer(valBytes)
	gobDecoder := gob.NewDecoder(gobBuffer)
	err = gobDecoder.Decode(result)
	return err
}

// DeleteObject removes an object from the database.
func (kv *MemoryStore) DeleteObject(table, id string) (err error) {
	key := [2]string{table, id}
	_, ok := kv.items[key]
	if !ok {
		return nil
	}
	kv.items[key] = nil
	return nil
}
