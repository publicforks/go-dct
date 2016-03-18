Docker Compose Test environment for go tests
==============================================
[![GoDoc](https://godoc.org/github.com/kbudde/go-dct?status.svg)](https://godoc.org/github.com/kbudde/go-dct)

Go-dct leverages ```go test``` to prepare integration test environments on demand for running tests against real services use docker compose. 

## features

- create docker containers using docker-compose for describing services
- dynamic ports allowed. No need to set fixed port in compose file. dct will get the used port from compose.
- useable from test methods
- support remote docker hosts ( e.g. windows/osx with docker-machine)

## planned


- support scaling/ serviceByID ?
- create helper around common services (e.g. mongoDB: backup,restore, count,...) in a separate project 
- Improve godoc
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
    
    //asking for host:port for the service with port 27017. Host will be 127.0.0.1 or ip from docker-machine remote docker_host
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

## Note on docker-machine

You must use docker-machine to set environment variables before running the tests on windows (and probably on osx as well).

Powershell: & "C:\Program Files\Docker Toolbox\docker-machine.exe" env default | Invoke-Expression

Cmd: FOR /f "tokens=*" %i IN ('docker-machine env default') DO %i