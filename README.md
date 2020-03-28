# kvstore

[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/flick-web/kvstore)

kvstore is a library that allows programs to have some sort of database-like persistence with very little setup. It provides an interface, `KeyValueStore`, which represents some memory or file-backed storage. Only `Get`, `Set`, and `Delete` methods are required.

kvstore also provides implementations of `KeyValueStore`:

- `SqliteStore`: Created with `NewSqliteDB(filename string)`. A key-value store backed by a sqlite3 database. Best for small projects where an extra database instance isn't worth the effort.
- `MemoryStore`: Created with `NewMemoryStore()`. An in-memory key-value store without persistence. Useful for testing.
