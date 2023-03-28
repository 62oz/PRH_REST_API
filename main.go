package main

import "fmt"

func main() {
	postalCodes := []string{"02100", "00140", "00930", "00710", "01730", "00500", "01760", "01690", "00510", "00180"}
	nCompanies := 10
	// Get the data from the API for each postal code
	data := make([]interface{}, len(postalCodes))
	for i, postalCode := range postalCodes {
		data[i] = GetHandler(postalCode, nCompanies)
	}

	fmt.Println(data[0])
}
