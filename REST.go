package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

func handlePostalCodeCompanies(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(companies)
}

func getCompaniesByPostalCode(db *sql.DB, code string) ([]Company, error) {
	rows, err := db.Query("SELECT * FROM companies WHERE postal_code = ?", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var company Company
		err := rows.Scan(&company.BusinessID, &company.Name, &company.RegistrationDate, &company.CompanyForm, &company.PostalCode, &company.DetailsURI)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}
