package repository

import (
	"fmt"
	"recipt-processor/model"

	"github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

func SetDB(database *memdb.MemDB) {
	db = database
}

func InsertTransaction(transaction *model.Transaction) error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	txn := db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("transaction", transaction); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func GetTransaction(uuid string) (*model.Transaction, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("transaction", "id", uuid)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	transaction := raw.(*model.Transaction)
	return transaction, nil
}

func GetAllTransactions() ([]*model.Transaction, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	txn := db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("transaction", "id")
	if err != nil {
		return nil, err
	}

	var transactions []*model.Transaction
	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		transaction := raw.(*model.Transaction)
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func DeleteSpecificTransaction(uuid string) error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	txn := db.Txn(true)
	defer txn.Abort()

	err := txn.Delete("transaction", &model.Transaction{Uuid: uuid})
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func DeleteAllTransactions() error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	txn := db.Txn(true)
	defer txn.Abort()

	iter, err := txn.Get("transaction", "id")
	if err != nil {
		return err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		transaction := raw.(*model.Transaction)
		if err := txn.Delete("transaction", transaction); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}
