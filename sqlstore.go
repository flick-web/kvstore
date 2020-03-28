package kvstore

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	// Import sqlite3 database driver
	_ "github.com/mattn/go-sqlite3"
)

// SQLStore is an object similar to sql.DB that provides simple methods for create,
// read, update, and delete functionality on key-value items.
type SQLStore struct {
	db *sql.DB
}

// NewSqliteDB initializes a new SQLStore database connection.
func NewSqliteDB(name string) (kv *SQLStore, err error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS kv(
		table_name TEXT NOT NULL,
		id TEXT NOT NULL,
		val BLOB,
		PRIMARY KEY (table_name, id)
	);`)
	if err != nil {
		return nil, err
	}

	kv = &SQLStore{db}
	return kv, nil
}

// SetObject creates or updates the key-value pair.
func (kv *SQLStore) SetObject(table, id string, value interface{}) error {
	gobBuffer := new(bytes.Buffer)
	gobEncoder := gob.NewEncoder(gobBuffer)
	err := gobEncoder.Encode(value)
	if err != nil {
		return err
	}

	_, err = kv.db.Exec(
		"INSERT OR REPLACE INTO kv (table_name, id, val) VALUES(?, ?, ?);",
		table, id, gobBuffer.Bytes(),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetObject retrieves and decodes the stored value into result.
func (kv *SQLStore) GetObject(table, id string, result interface{}) (err error) {
	row := kv.db.QueryRow(
		"SELECT val FROM kv WHERE table_name = ? AND id = ?",
		table, id,
	)
	var buf []byte
	err = row.Scan(&buf)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrKeyNotFound
		}
		return err
	}

	gobBuffer := bytes.NewBuffer(buf)
	gobDecoder := gob.NewDecoder(gobBuffer)
	err = gobDecoder.Decode(result)
	return err
}

// DeleteObject removes an object from the database.
func (kv *SQLStore) DeleteObject(table, id string) (err error) {
	_, err = kv.db.Exec(
		"DELETE FROM kv WHERE table_name = ? AND id = ?",
		table, id,
	)
	return err
}
