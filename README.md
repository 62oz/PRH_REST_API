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
