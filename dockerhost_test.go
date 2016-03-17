package dct

import (
	"os"
	"testing"
)

func TestGetDockerHost(t *testing.T) {
	old := os.Getenv(dockerhostEnvvar)
	defer os.Setenv(dockerhostEnvvar, old)
	os.Unsetenv(dockerhostEnvvar)
	if host := getDockerHost(); host != "127.0.0.1" {
		t.Errorf("If %v is not set 127.0.0.1 should used as docker host. Found: %v", dockerhostEnvvar, host)
	}

	os.Setenv(dockerhostEnvvar, "tcp://192.168.99.100:2376")
	if host := getDockerHost(); host != "192.168.99.100" {
		t.Errorf("If %v is set the ip of the host should be used as docker host. Found: %v", dockerhostEnvvar, host)
	}

	os.Setenv(dockerhostEnvvar, "tcp://docker-host:2376")
	if host := getDockerHost(); host != "docker-host" {
		t.Errorf("If %v is set to a name this should be used as docker host. Found: %v", dockerhostEnvvar, host)
	}

	os.Setenv(dockerhostEnvvar, "tcp://docker-host.mydomain.local:2376")
	if host := getDockerHost(); host != "docker-host.mydomain.local" {
		t.Errorf("If %v is set to a fqdn this should be used as docker host. Found: %v", dockerhostEnvvar, host)
	}

	os.Setenv(dockerhostEnvvar, "tcp://[::1]:2376")
	if host := getDockerHost(); host != "[::1]" {
		t.Errorf("If %v is set the ip of the host should be used as docker host for ipv6 as well. Found: %v", dockerhostEnvvar, host)
	}

	os.Setenv(dockerhostEnvvar, "tcp://[fe80::a00:27ff:fefa:c5cb]:2376")
	if host := getDockerHost(); host != "[fe80::a00:27ff:fefa:c5cb]" {
		t.Errorf("If %v is set the ip of the host should be used as docker host for ipv6 as well. Found: %v", dockerhostEnvvar, host)
	}
}
