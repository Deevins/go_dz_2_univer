package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Transaction struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func main() {
	r := mux.NewRouter()

	// /ping для проверки работы
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}).Methods(http.MethodGet)

	// /transactions проксирует запрос к Ledger
	r.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:9090/transactions")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to reach ledger: %v", err), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "ledger returned non-200 status", http.StatusBadGateway)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read ledger response", http.StatusInternalServerError)
			return
		}

		var transactions []Transaction
		if err := json.Unmarshal(body, &transactions); err != nil {
			http.Error(w, "failed to parse ledger response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
	}).Methods(http.MethodGet)

	port := 8080
	fmt.Printf("Gateway service started on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
