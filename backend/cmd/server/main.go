package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eli-bosch/DBMS-final/config"
	_ "github.com/eli-bosch/DBMS-final/internal/db"
	"github.com/eli-bosch/DBMS-final/internal/routes"
	"github.com/gorilla/mux"
)

func main() {

	db := config.Connect()
	defer db.Close()

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	routes.DBMSRoutes(r)

	http.Handle("/", r)
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe("localhost:9010", r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://my-static-site-eli.s3-website.us-east-2.amazonaws.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
