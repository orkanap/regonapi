# Go REGON API [![GoDoc](https://godoc.org/github.com/orkanap/regonapi?status.svg)](https://godoc.org/github.com/orkanap/regonapi)

This project is Go interface to Polish National Official Business Register (REGON)
database. More information [here](https://bip.stat.gov.pl/en/regon/). Library is a
thin wrapper for official SOAP webservice API. Check technical documentation
[here](https://api.stat.gov.pl/Home/RegonApi?lang=en).

## Features

* search for business entities by NIP (tax id), REGON (statistical numbers)
  or KRS (National Court Register)
* get system status
* get basic reports
* production and test (sandbox) environments
* implements BIR1 version 1.1
* context-aware HTTP requests

To access production environment you need user key issued by REGON
administrators.

## Installation

```go
go get github.com/orkanap/regonapi
```

## Basic usage

Find entities by NIP.

```go
// New client, empty key connects to test envinroment
regonsvc := regonapi.NewClient(context.Background(), os.Getenv("REGON_API_KEY"))

err := regonsvc.Login()
if err != nil {
    log.Fatal(err)
}
defer regonsvc.Logout()

entities, err := regonsvc.SearchByNIP("5261040828")
if err != nil {
    log.Fatal(err)
}

for _, entity := range entities {
    fmt.Println("[", entity.Type, "]")
    fmt.Println("NIP:", entity.NIP, "REGON:", entity.REGON)
    fmt.Println(entity.Name)
    fmt.Println(entity.Street, entity.PropertyNumber, entity.ApartmentNumber)
    fmt.Println(entity.PostalCode, entity.City)
}
```

## Get more details

```go
// List PKDs
pkds, err := regonsvc.LegalPersonPKDList("340771731")
if err != nil {
    log.Fatal(err)
}

fmt.Println("PKD:")
for _, pkd := range pkds {
    fmt.Println(pkd.Code, pkd.Name)
}
```

For complete example check `./example` folder,  `integration_test.go` source
file or `godoc`.

## Tests

Integration tests make http requests to sandbox endpoint. Beware of API quotas.
To perform tests `make test` or `make test-unit` for unit tests only.

## Limitations

* searching by group of identifiers is not implemented
* only basic reports are implemented
* no control of NIP, REGON, KRS identifiers

## License

The MIT License (MIT) - see [`LICENSE.md`](https://github.com/orkanap/regonapi/blob/master/LICENSE.md) for more details
