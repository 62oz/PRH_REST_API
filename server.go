package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlePostalCodeCompanies(w http.ResponseWriter, r *http.Request) {
	// Check for the correct path
	prefix := "/postal_codes/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, prefix)
	if !strings.HasSuffix(code, "/companies") {
		http.NotFound(w, r)
		return
	}
	code = strings.TrimSuffix(code, "/companies")

	db := getDB("companies.db")
	companies, err := getCompaniesByPostalCode(db, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If data is not in the database, fetch it from the API
	if len(companies) == 0 {
		// Get the data from the API
		log.Println("Fetching data from the API for new postal code...")
		companies = GetHandler(code, 2)
		if companies == nil {
			http.Error(w, "No companies found in the database for this postal code.", http.StatusNotFound)
			return
		}
		// Insert the data into the database
		err := insertCompanies(db, companies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(companies)
}
