package kvstore_test

import (
	"testing"

	"github.com/flick-web/kvstore"
)

type testObj struct {
	X map[string]string
	Y string
	Z int
}

func TestSqliteStore(t *testing.T) {
	memDB := kvstore.NewMemoryStore()
	sqliteDB, err := kvstore.NewSqliteDB("keyvalue.db")
	if err != nil {
		t.Error(err)
	}

	dbs := []kvstore.KeyValueStore{sqliteDB, memDB}
	for _, db := range dbs {
		testVal := testObj{
			X: map[string]string{"one": "1", "two": "2"},
			Y: "Hello!",
			Z: 42,
		}
		err = db.Set("test", "12345", testVal)
		if err != nil {
			t.Error(err)
		}

		out := testObj{}
		err = db.Get("test", "67890", &out)
		if err == nil {
			t.Error("Expected ErrKeyNotFound error")
		}
		if err != kvstore.ErrKeyNotFound {
			t.Error(err)
		}

		err = db.Get("test", "12345", &out)
		if err != nil {
			t.Error(err)
		}
		if out.Y != testVal.Y || out.Z != testVal.Z {
			t.Error("Objects did not match")
		}

		db.Delete("test", "12345")
	}
}
