package dct

import (
	"os"
	"testing"
)

func TestNewComposer_WithNoDockerCompose(t *testing.T) {
	path := os.Getenv("PATH")
	os.Unsetenv("PATH")
	_, err := NewComposer("")
	if err != ErrDockerComposeNotFound {
		t.Errorf("Path is empty. It should not be possible to find docker-compose cmd in path. Err:%v", err)
	}
	os.Setenv("PATH", path)
}

func TestInvalidComposeFile(t *testing.T) {
	c, err := NewComposer("testFiles/docker-compose-invalid.yml")
	if err != nil {
		t.Error(err)
	}
	err = c.StartAll()
	if err == nil {
		t.Errorf("Compose file is invalid. Error was inspected")
	}

	count, err := c.ServiceCount()
	if count > -1 {
		t.Error("File is invalid. ServiceCount is wrong")
	}
	if err == nil {
		t.Errorf("Error was expected. Err:%v", err)
	}
}

func TestRemoveAll(t *testing.T) {
	c, err := NewComposer("testFiles/docker-compose.yml")
	if err != nil {
		t.Error(err)
	}

	//RemoveAll container and check success
	err = c.RemoveAll()
	if err != nil {
		t.Error(err)
	}
	count, err := c.ServiceCount()
	if err != nil {
		t.Error(err)
	}
	if count > 0 {
		t.Errorf("All images should have been removed. Found: %v", count)
	}

	//StartAll, StopAll -> Count should be 2
	err = c.StartAll()
	if err != nil {
		t.Error(err)
	}
	err = c.StopAll()
	if err != nil {
		t.Error(err)
	}
	count, err = c.ServiceCount()
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("2 Images expected. Found: %v", count)
	}

	//RemoveAll, Start One, count -> 1
	err = c.RemoveAll()
	if err != nil {
		t.Error(err)
	}
	err = c.Service("nginx").Start()
	if err != nil {
		t.Error(err)
	}
	count, err = c.ServiceCount()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("1 Image expected. Found: %v", count)
	}

	//RemoveAll (1 running). count -> 0
	err = c.RemoveAll()
	if err != nil {
		t.Error(err)
	}
	count, err = c.ServiceCount()
	if err != nil {
		t.Error(err)
	}
	if count > 0 {
		t.Errorf("All images should have been removed. Found: %v", count)
	}

	//RemoveAll (0 running). count -> 0
	err = c.RemoveAll()
	if err != nil {
		t.Error(err)
	}
	count, err = c.ServiceCount()
	if err != nil {
		t.Error(err)
	}
	if count > 0 {
		t.Errorf("All images should have been removed. Found: %v", count)
	}
}
