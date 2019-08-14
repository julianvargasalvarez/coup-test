package main

import (
	"app/controllers/scooters"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Running server...")
	http.HandleFunc("/scooters/", scooters.Index)
	http.ListenAndServe(":8080", nil)
}
