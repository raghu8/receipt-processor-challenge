package test

import (
	"fmt"
	"recipt-processor/model"
	"recipt-processor/repository"
	"recipt-processor/service"
	"testing"

	"github.com/hashicorp/go-memdb"
)

func initMockDB() {
	schema := &memdb.DBSchema{
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

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		fmt.Println("Failed to initialize mock DB:", err)
	}
	repository.SetDB(db)
}

func TestCalculatePoints(t *testing.T) {
	initMockDB()

	tests := []struct {
		name        string
		transaction model.Transaction
		expected    int
	}{
		{
			name: "Alphanumeric characters in retailer name",
			transaction: model.Transaction{
				Retailer:     "M&M Corner Market",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "1.01",
			},
			expected: 14,
		},
		{
			name: "Round dollar amount",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "43.00",
			},
			expected: 81, // 50 points for round dollar amount
		},
		{
			name: "Total is a multiple of 0.25",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "25.25",
			},
			expected: 31,
		},
		{
			name: "5 points for every two items",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items: []model.Item{
					{ShortDescription: "Item 1", Price: "1.00"},
					{ShortDescription: "Item 2", Price: "1.00"},
				},
				Total: "2.41",
			},
			expected: 13,
		},
		{
			name: "Item description length multiple of 3",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items: []model.Item{
					{ShortDescription: "ABC", Price: "5.00"},
				},
				Total: "5.98",
			},
			expected: 7, //
		},
		{
			name: "Odd day in purchase date",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "5.34",
			},
			expected: 12,
		},
		{
			name: "Time of purchase after 2:00pm and before 4:00pm",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "15:00",
				Items:        []model.Item{},
				Total:        "4.51",
			},
			expected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateTransaction(&tt.transaction)
			if err != nil {
				t.Fatalf("CreateTransaction() error = %sv", err)
			}

			pointsModel, err := service.GetTransaction(tt.transaction.Uuid)
			if err != nil {
				t.Fatalf("CalculatePoints() error = %v", err)
			}
			if pointsModel.Points != tt.expected {
				t.Errorf("CalculatePoints() = %v, want %v", pointsModel.Points, tt.expected)
			}
		})
	}
}

func TestGetAllTransactions(t *testing.T) {
	initMockDB()

	tests := []struct {
		name        string
		transaction model.Transaction
	}{
		{
			name: "Get all transactions",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "1.01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateTransaction(&tt.transaction)
			if err != nil {
				t.Fatalf("CreateTransaction() error = %v", err)
			}

			transactions, err := service.GetAllTransactions()
			if err != nil {
				t.Fatalf("GetAllTransactions() error = %v", err)
			}
			if len(transactions) != 1 {
				t.Errorf("GetAllTransactions() = %v, want %v", len(transactions), 1)
			}
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
	initMockDB()

	tests := []struct {
		name        string
		transaction model.Transaction
	}{
		{
			name: "Delete transaction",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "1.01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateTransaction(&tt.transaction)
			if err != nil {
				t.Fatalf("CreateTransaction() error = %v", err)
			}

			err = service.DeleteTransaction(tt.transaction.Uuid)
			if err != nil {
				t.Fatalf("DeleteTransaction() error = %v", err)
			}

			transactions, err := service.GetAllTransactions()
			if err != nil {
				t.Fatalf("GetAllTransactions() error = %v", err)
			}
			if len(transactions) != 0 {
				t.Errorf("GetAllTransactions() = %v, want %v", len(transactions), 0)
			}
		})
	}
}

func TestDeleteSpecifiedTransaction(t *testing.T) {
	initMockDB()

	tests := []struct {
		name        string
		transaction model.Transaction
	}{
		{
			name: "Delete specified transaction",
			transaction: model.Transaction{
				Retailer:     "Target",
				Uuid:         "12345",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:01",
				Items:        []model.Item{},
				Total:        "1.01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateTransaction(&tt.transaction)
			if err != nil {
				t.Fatalf("CreateTransaction() error = %v", err)
			}

			err = service.DeleteSpecifiedTransaction(tt.transaction.Uuid)
			if err != nil {
				t.Fatalf("DeleteSpecifiedTransaction() error = %v", err)
			}

			transactions, err := service.GetAllTransactions()
			if err != nil {
				t.Fatalf("GetAllTransactions() error = %v", err)
			}
			if len(transactions) != 0 {
				t.Errorf("GetAllTransactions() = %v, want %v", len(transactions), 0)
			}
		})
	}
}
