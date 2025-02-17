
# Go-Plate
A Golang API boilerplate

## Installation

The project use Fiber which needs golang >1.16 to be installed on your computer. Then install the project in your working directory

```bash
$ git clone https://github.com/medrupaloscil/go-plate
$ cd go-plate
$ go get .
$ cp .env.example .env
$ go run main.go
```
The database is empty by default, you can importe the fixtures/fixtures.sql file in your own database to have needed data

## Project structure

The project is structured as following :

```
| - controller
| - logs
| - models
| - routing
| - services
| - tests
| - translations
```

### Controllers

This folder contains the controllers, which means all the logic and data treatment before storing in database

### Controllers
Containing the logs. When the log file is more than 10Mb, it creates a new file.

### Models

This is where models are structured and complex database function are developed.
Model files are in kebab-case because all *_test.go files are considered test files and won't be interpreted in dev/prod environment.

### Routing

This is where routes are declared. Every subroute has its own file to properly separate API parts.

### Services

This folder's purpose is to centralize all important/big reusable funtions (as token manager, data validators, s3 storage, ...)

### Tests

This folder contains all the functional tests

### Translations

This folder contains all the lang files in yaml format

## Automated tests

To run the tests :

```bash
$ go get github.com/stretchr/testify/assert
$ go test -v ./tests
```