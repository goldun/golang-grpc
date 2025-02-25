# quotes-service

Basic Golang service with GRPC server and corresponding client code

## Features
Returns a quote(s) for a specified source.
Source could be anything: person, character, film, book. etc
If source is missing, or Data Store ain't have any quotes for a given source an error will be returned

### Data Storage
Update `testdata/quotes.json` to add more data

## How to run

### Server
```shell
go run server/server.go
```

### Client
```shell
go run client/client.go
```