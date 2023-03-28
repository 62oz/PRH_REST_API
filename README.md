# Companies by postal code

This project fetches data about companies in the Helsinki metropolitan area registered to a given postal code using the [Finnish Business Information System (BIS) API](https://avoindata.prh.fi/ytj_en.html). It then stores the data in a SQLite database and provides a REST API to retrieve the data in JSON format.

## Requirements

- Go 1.16 or higher
- SQLite 3

## Installation

Clone the repository:

```
git clone https://github.com/62oz/PRH_REST_API.git
```

Navigate to the project directory:

```
cd  PRH_REST_API
```

Install the required Go packages:

```
go mod download
```

## Usage
To start the API server, run:

```
go run .
```

Once the server is running, you can access the API by making GET requests to the following endpoint (for example you can go to this url directly in your browser):

```
http://localhost:8080/postal_codes/[POSTAL_CODE]/companies
```
Replace [POSTAL_CODE] with the desired postal code. For example:

```
http://localhost:8080/postal_codes/00100/companies
```
This will return a JSON array of companies registered to the postal code "00100".

## Notes
The REST API will by default get and store data of 20 companies from the following postal codes: {"02100", "00140", "00930", "00710", "01730", "00500", "01760", "01690", "00510", "00180"}. In addition, it will also fetch the data of 20 companies from any postal code specified in your GET request, if it exists.

If you want to modify the number of companies fetched for a postal code, you can do so by modifying the value of nCompanies in the main function (main.go):
```
func main() {
	// Choose the postal codes and the number of companies to get for each
	postalCodes := []string{"02100", "00140", "00930", "00710", "01730", "00500", "01760", "01690", "00510", "00180"}
	nCompanies := 20
  
  // Rest of the code ...
```
