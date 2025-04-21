package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eli-bosch/DBMS-final/internal/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.DBMSRoutes(r)

	http.Handle("/", r)
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe("localhost:9010", r))

}
