package kvstore

import "testing"

type testObj struct {
	X map[string]string
	Y string
	Z int
}

func TestSqliteStore(t *testing.T) {
	memDB := NewMemoryStore()
	sqliteDB, err := NewSqliteDB("keyvalue.db")
	if err != nil {
		t.Error(err)
	}

	dbs := []KeyValueStore{sqliteDB, memDB}
	for _, db := range dbs {
		testVal := testObj{
			X: map[string]string{"one": "1", "two": "2"},
			Y: "Hello!",
			Z: 42,
		}
		err = db.SetObject("test", "12345", testVal)
		if err != nil {
			t.Error(err)
		}

		out := testObj{}
		err = db.GetObject("test", "67890", &out)
		if err == nil {
			t.Error("Expected ErrKeyNotFound error")
		}
		if err != ErrKeyNotFound {
			t.Error(err)
		}

		err = db.GetObject("test", "12345", &out)
		if err != nil {
			t.Error(err)
		}
		if out.Y != testVal.Y || out.Z != testVal.Z {
			t.Error("Objects did not match")
		}

		db.DeleteObject("test", "12345")
	}
}
