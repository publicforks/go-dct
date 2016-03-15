package dct

import (
	"testing"
)

func TestPort(t *testing.T) {
	comp, _ := NewComposer("testFiles/docker-compose.yml")
	err := comp.StartAll()
	defer comp.StopAll()
	if err != nil {
		t.Errorf("Error starting services: %v", err)
	}

	if res, err := comp.Service("nginx").Port("80"); err != nil {
		t.Errorf("Error getting service: %v", err)
	} else if res != "8088" {
		t.Errorf("Wrong Port for 8088: %v", res)
	}

	if res, err := comp.Service("nginx").Port("443"); err != nil {
		t.Errorf("Error getting service port: %v", err)
	} else if res != "44443" {
		t.Errorf("Wrong Port for 44443: %v", res)
	}

	if res, err := comp.Service("nginx").Port("11111"); err != ErrUnknownServiceOrPort {
		t.Errorf("UnexpectedError: %v", err)
	} else if res != "" {
		t.Errorf("Wrong Port: %v", res)
	}
}

func TestURI(t *testing.T) {
	comp, _ := NewComposer("testFiles/docker-compose.yml")
	err := comp.StartAll()
	defer comp.StopAll()
	if err != nil {
		t.Errorf("Error starting services: %v", err)
	}

	if res, err := comp.Service("nginx").URI("80"); err != nil {
		t.Errorf("Error getting service: %v", err)
	} else if res != "127.0.0.1:8088" {
		t.Errorf("Wrong URI for 127.0.0.1:80: %v", res)
	}

	if res, err := comp.Service("nginx").URI("443"); err != nil {
		t.Errorf("Error getting service: %v", err)
	} else if res != "127.0.0.1:44443" {
		t.Errorf("Wrong URI for 127.0.0.1:443: %v", res)
	}

	if res, err := comp.Service("nginx").URI("11111"); err != ErrUnknownServiceOrPort {
		t.Errorf("UnexpectedError: %v", err)
	} else if res != "" {
		t.Errorf("Wrong Port: %v", res)
	}
}

func TestStartService(t *testing.T) {
	comp, err := NewComposer("testFiles/docker-compose.yml")
	if err != nil {
		t.Error(err)
		return
	}
	nginx := comp.Service("nginx")
	err = nginx.Stop()
	if err != nil {
		t.Error("Stop should work for an existing service")
	}
	p, err := nginx.Port("80")
	if err == nil {
		t.Error("Service is stopped. Error expected.")
	}
	if p != "" {
		t.Errorf("Service is stopped. Empty port expected.Found: %v", p)
	}

	err = nginx.Start()
	if err != nil {
		t.Error("Start should work for this service")
	}
	p, err = nginx.Port("80")
	if err != nil {
		t.Errorf("Service is started. Port should return mapping. Got error:%v", err)
	}
	if p != "8088" {
		t.Errorf("Service is started and mapped to 8088. Result was %v", p)
	}

	err = comp.Service("doesNotExist").Start()
	if err == nil {
		t.Errorf("Service does not exist -> it should not start")
	}
}
