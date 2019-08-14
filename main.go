package main

import (
        "fmt"
	"net/http"
	"app/controllers/scooters"
)

func main() {
        fmt.Println("Running server...")
	http.HandleFunc("/scooters/", scooters.Index)
	http.ListenAndServe(":8080", nil)
}
