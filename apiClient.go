package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Company struct {
	BusinessID       string `json:"businessId"`
	Name             string `json:"name"`
	RegistrationDate string `json:"registrationDate"`
	CompanyForm      string `json:"companyForm"`
	PostalCode       string
	DetailsURI       string `json:"detailsUri"`
}

func GetHandler(postalCode string, nCompanies int) []Company {
	// Get the data from the API
	url := "https://avoindata.prh.fi/bis/v1?totalResults=false&maxResults=" + strconv.Itoa(nCompanies) + "&resultsFrom=0&streetAddressPostCode=" + postalCode
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Get request error:", err)
		return nil
	}
	defer resp.Body.Close()

	// Check for bad status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status Error:", resp.StatusCode)
		return nil
	}

	var result struct {
		Results []Company `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Decoder error:", err)
		return nil
	}

	// Add postal code to each company
	for i := range result.Results {
		result.Results[i].PostalCode = postalCode
	}

	return result.Results
}
