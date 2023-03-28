package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func getDB(path string) *sql.DB {
	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	// Create the companies table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS companies (
        business_id TEXT PRIMARY KEY UNIQUE,
        name TEXT NOT NULL,
        registration_date TEXT NOT NULL,
        company_form TEXT NOT NULL,
		postal_code TEXT NOT NULL,
		details_uri TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertCompanies(db *sql.DB, companies []Company) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO companies (business_id, name, registration_date, company_form, postal_code, details_uri) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, company := range companies {
		_, err = stmt.Exec(company.BusinessID, company.Name, company.RegistrationDate, company.CompanyForm, company.PostalCode, company.DetailsURI)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
