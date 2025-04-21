package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
    Content string `json:"content"`
}

func main() {
    db, err := sql.Open("mysql", "casonp:aC5eequa@tcp(localhost:3306)/testdb")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    http.HandleFunc("api/message", func(w http.ResponseWriter, r *http.Request) {
        var msg Message
        err := db.QueryRow("SELECT content FROM messages LIMIT 1").Scan(&msg.Content)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(msg)
    })

	fmt.Print("Listening on port :8080")
    http.ListenAndServe(":8080", nil)
}