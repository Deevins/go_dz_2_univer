package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Transaction — модель финансовой транзакции.
type Transaction struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

var transactions []Transaction

// AddTransaction добавляет новую транзакцию.
func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("transaction amount cannot be zero")
	}
	tx.ID = len(transactions) + 1
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}
	transactions = append(transactions, tx)
	return nil
}

// ListTransactions возвращает все транзакции.
func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

// HTTP-обработчик для /transactions
func handleListTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListTransactions())
}

func main() {
	fmt.Println("Ledger service started on port 9090")

	_ = AddTransaction(Transaction{Amount: 120.50, Category: "Food", Description: "Dinner"})
	_ = AddTransaction(Transaction{Amount: 300, Category: "Transport", Description: "Taxi ride"})
	_ = AddTransaction(Transaction{Amount: 750, Category: "Entertainment", Description: "Concert tickets"})

	r := mux.NewRouter()
	r.HandleFunc("/transactions", handleListTransactions).Methods(http.MethodGet)

	http.ListenAndServe(":9090", r)
}
