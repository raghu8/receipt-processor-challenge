package controller

import (
	"encoding/json"
	"net/http"
	"recipt-processor/model"
	"recipt-processor/repository"
	"recipt-processor/service"

	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(r *mux.Router) {
	r.HandleFunc("/transactions", CreateTransaction).Methods("POST")
	r.HandleFunc("/transactions/{uuid}", GetTransaction).Methods("GET")
	r.HandleFunc("/transactions", GetAllTransactions).Methods("GET")
	r.HandleFunc("/delete", DeleteAllTransactions).Methods("DELETE")
	r.HandleFunc("/delete/{uuid}", DeleteSpecificTransaction).Methods("DELETE")
	r.HandleFunc("/health", HealthCheck).Methods("GET")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction model.Transaction
	var transactionInfo model.TransactionInfo
	newUUID := uuid.New()
	transaction.Uuid = newUUID.String()
	fmt.Print("UUID: ", transaction.Uuid)
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		fmt.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.InsertTransaction(&transaction); err != nil {
		fmt.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transactionInfo.ID = newUUID.String()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transactionInfo)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuid := params["uuid"]

	transaction, err := service.GetTransaction(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if transaction == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := service.GetAllTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func DeleteAllTransactions(w http.ResponseWriter, r *http.Request) {
	err := repository.DeleteAllTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteSpecificTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuid := params["uuid"]

	err := service.DeleteTransaction(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
