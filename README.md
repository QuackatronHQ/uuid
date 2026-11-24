# uuid
The uuid package generates and inspects UUIDs based on
[RFC 9562](https://datatracker.ietf.org/doc/html/rfc9562)
and DCE 1.1: Authentication and Security Services. 

This package is based on the github.com/pborman/uuid package (previously named
code.google.com/p/go-uuid).  It differs from these earlier packages in that
a UUID is a 16 byte array rather than a byte slice.  One loss due to this
change is the ability to represent an invalid UUID (vs a NIL UUID).

###### Install
```sh
go get github.com/google/uuid
```

###### Usage
```go
id := uuid.New()                 // random, Version 4
sortable, err := uuid.NewV7()    // time ordered, Version 7
backfill, err := uuid.NewV7WithTime(createdAt)

parsed, err := uuid.Parse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
if errors.Is(err, uuid.ErrInvalidLength) {
	// handle bad input
}

if ts, ok := sortable.Timestamp(); ok {
	fmt.Println("created at", ts)
}
```

Go 1.19 or later is required.

###### Choosing a version
| Version | Constructor | Use when |
|---------|-------------|----------|
| 4 | `New`, `NewRandom` | You need an opaque identifier |
| 7 | `NewV7`, `NewV7WithTime` | Identifiers should sort by creation time, such as database keys |
| 5 | `NewSHA1` | The same name must always map to the same UUID |
| 8 | `NewV8` | You need to pack your own bits into a UUID |
| 1, 6 | `NewUUID`, `NewV6` | You must interoperate with existing time and node based UUIDs |

###### Documentation 
[![Go Reference](https://pkg.go.dev/badge/github.com/google/uuid.svg)](https://pkg.go.dev/github.com/google/uuid)

Full `go doc` style documentation for the package can be viewed online without
installing this package by using the GoDoc site here: 
https://pkg.go.dev/github.com/google/uuid
