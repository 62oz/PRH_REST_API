package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestInsertCompanies(t *testing.T) {
	// Open a test database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create the companies table
	_, err = db.Exec("CREATE TABLE companies (business_id TEXT, name TEXT, registration_date TEXT, company_form TEXT, postal_code TEXT, details_uri TEXT)")
	if err != nil {
		t.Fatal(err)
	}

	// Insert some test data
	companies := []Company{
		{BusinessID: "123", Name: "Company A", RegistrationDate: "2021-01-01", CompanyForm: "Form A", PostalCode: "00100", DetailsURI: "https://example.com/a"},
		{BusinessID: "456", Name: "Company B", RegistrationDate: "2021-02-01", CompanyForm: "Form B", PostalCode: "00100", DetailsURI: "https://example.com/b"},
		{BusinessID: "789", Name: "Company C", RegistrationDate: "2021-03-01", CompanyForm: "Form C", PostalCode: "00200", DetailsURI: "https://example.com/c"},
	}
	err = insertCompanies(db, companies)
	if err != nil {
		t.Fatal(err)
	}

	// Query the data to check if it was inserted correctly
	rows, err := db.Query("SELECT business_id, name, registration_date, company_form, postal_code, details_uri FROM companies")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	// Check if the data matches the inserted data
	var results []Company
	for rows.Next() {
		var c Company
		err := rows.Scan(&c.BusinessID, &c.Name, &c.RegistrationDate, &c.CompanyForm, &c.PostalCode, &c.DetailsURI)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, c)
	}
	if !reflect.DeepEqual(companies, results) {
		t.Errorf("Expected %v, got %v", companies, results)
	}
}

func testGetDB(t *testing.T) {
	// Create temporary file for database
	tmpfile, err := ioutil.TempFile("", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Test opening database
	db := getDB(tmpfile.Name())
	if db == nil {
		t.Error("Database not returned")
	}

	// Test creating table
	table := "companies"
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", table)
	rows, err := db.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Errorf("Table %s not created", table)
	}
}

func TestHandlePostalCodeCompanies(t *testing.T) {
	// Set up mock database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening mock database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE companies (
		business_id TEXT PRIMARY KEY,
		name TEXT,
		registration_date TEXT,
		company_form TEXT,
		postal_code TEXT,
		details_uri TEXT
	)`)
	if err != nil {
		t.Fatalf("Error creating mock companies table: %v", err)
	}

	// Insert test data into database
	companies := []Company{
		{
			BusinessID:       "3353862-4",
			Name:             "Stay Sharp Oy",
			RegistrationDate: "2023-03-23",
			CompanyForm:      "OY",
			DetailsURI:       "http://avoindata.prh.fi/opendata/bis/v1/3353862-4",
			PostalCode:       "00210",
		},
		{
			BusinessID:       "3349694-6",
			Name:             "HSN Consulting Oy",
			RegistrationDate: "2023-03-06",
			CompanyForm:      "OY",
			DetailsURI:       "http://avoindata.prh.fi/opendata/bis/v1/3349694-6",
			PostalCode:       "00210",
		},
	}
	err = insertCompanies(db, companies)
	if err != nil {
		t.Fatalf("Error inserting test data into mock database: %v", err)
	}

	// Set up mock HTTP request/response objects
	req := httptest.NewRequest("GET", "/postal_codes/00210/companies", nil)
	w := httptest.NewRecorder()

	// Call the handler function
	handlePostalCodeCompanies(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// Decode the response body
	var respBody []Company
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check that the response body matches the test data
	if !reflect.DeepEqual(respBody, companies) {
		t.Errorf("Expected response body %+v, but got %+v", companies, respBody)
	}
}

func TestHandlePostalCodeCompaniesNotFound(t *testing.T) {
	// Set up mock database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening mock database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE companies (
		business_id TEXT PRIMARY KEY,
		name TEXT,
		registration_date TEXT,
		company_form TEXT,
		postal_code TEXT,
		details_uri TEXT
	)`)
	if err != nil {
		t.Fatalf("Error creating mock companies table: %v", err)
	}

	// Set up mock HTTP request/response objects
	req := httptest.NewRequest("GET", "/postal_codes/12345/companies", nil)
	w := httptest.NewRecorder()

	// Call the handler function
	handlePostalCodeCompanies(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestGetHandler(t *testing.T) {
	// Excpeted response body
	expected := []Company{
		{
			BusinessID:       "3353862-4",
			Name:             "Stay Sharp Oy",
			RegistrationDate: "2023-03-23",
			CompanyForm:      "OY",
			DetailsURI:       "http://avoindata.prh.fi/opendata/bis/v1/3353862-4",
			PostalCode:       "00210",
		},
		{
			BusinessID:       "3349694-6",
			Name:             "HSN Consulting Oy",
			RegistrationDate: "2023-03-06",
			CompanyForm:      "OY",
			DetailsURI:       "http://avoindata.prh.fi/opendata/bis/v1/3349694-6",
			PostalCode:       "00210",
		},
	}

	// Call the handler function
	result := GetHandler("00210", 2)

	// Check that the response body matches the test data
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected response body %+v, but got %+v", expected, result)
	}
}
