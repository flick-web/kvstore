package kvstore

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	// Import sqlite3 database driver
	_ "github.com/mattn/go-sqlite3"
)

// SqliteStore is an object similar to sql.DB that provides simple methods for create,
// read, update, and delete functionality on key-value items.
type SqliteStore struct {
	db *sql.DB
}

// NewSqliteDB initializes a new SqliteStore database connection.
func NewSqliteDB(name string) (kv *SqliteStore, err error) {
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

	kv = &SqliteStore{db}
	return kv, nil
}

// Set creates or updates the key-value pair.
func (kv *SqliteStore) Set(table, id string, value interface{}) error {
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

// Get retrieves and decodes the stored value into result.
func (kv *SqliteStore) Get(table, id string, result interface{}) (err error) {
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

// Delete removes an object from the database.
func (kv *SqliteStore) Delete(table, id string) (err error) {
	_, err = kv.db.Exec(
		"DELETE FROM kv WHERE table_name = ? AND id = ?",
		table, id,
	)
	return err
}
