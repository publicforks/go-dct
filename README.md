Docker Compose Test environment for go tests
==============================================
[![GoDoc](https://godoc.org/github.com/kbudde/go-dct?status.svg)](https://godoc.org/github.com/kbudde/go-dct)

Go-dct leverages ```go test``` to prepare integration test environments on demand for running tests against real services use docker compose. 

## features

- create docker containers using docker-compose for describing services
- useable from test methods

## planned

- support remote docker hosts ( e.g. windows/osx with docker-machine)
- support scaling/ serviceByID ?
- create helper around common services (e.g. mongoDB: backup,restore, count,...) in a separate project 
- Improve iodoc
- generate projectnames for each test
- support docker-compose variable substitution

## Quickstart

```go
package integrationTest

import (
    "testing"
    "github.com/kbudde/go-dct"
)

func TestSave(t *testing.T) {
    testEnv, _ := dct.NewComposer("docker-compose.yml")
    testEnv.RemoveAll() //cleanup
    testEnv.StartAll()
    defer testEnv.RemoveAll()
    uri, _ := testEnv.Service("mongo").URI("27017")
    mongoURI := fmt.Sprintf("mongodb://%v/MyDB", uri)
    
    provider, _ := PizzaProvider(mongoURI)
    
    if provider.PizzaCount()!=0 {
        t.Error("No Pizza expected")
    }
    
    if err:= provider.NewPizza("tonno"); err!=nil {
        //...
    }
}
```

## production ready?

Not really. Useable perhaps.

It is the first release; so it is likely that the api will change. 
