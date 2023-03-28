package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Choose the postal codes and the number of companies to get for each
	postalCodes := []string{"02100", "00140", "00930", "00710", "01730", "00500", "01760", "01690", "00510", "00180"}
	nCompanies := 20

	log.Println("Fetching data from the API...")
	// Get the data from the API
	data := make([]interface{}, len(postalCodes))
	for i, postalCode := range postalCodes {
		data[i] = GetHandler(postalCode, nCompanies)
	}
	// Create the database
	db := getDB("companies.db")

	log.Println("Inserting data into the database...")
	// Insert the data into the database
	for _, postalCodes := range data {
		err := insertCompanies(db, postalCodes.([]Company))
		if err != nil {
			fmt.Println("Insert error:", err)
		}
	}

	// Start the server
	http.HandleFunc("/postal_codes/", handlePostalCodeCompanies)
	log.Println("Starting server port 8080 (http://localhost:8080/)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
