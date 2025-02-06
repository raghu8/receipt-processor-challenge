package configuration

import (
	"github.com/hashicorp/go-memdb"
)

var DB *memdb.MemDB

var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"transaction": {
			Name: "transaction",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Uuid"},
				},
			},
		},
	},
}

func InitDB() (*memdb.MemDB, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	return db, nil
}
