package dct

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

const dockerhostEnvvar = "DOCKER_HOST"

func getDockerHost() string {
	defaultHost := "127.0.0.1"

	dockerHost := os.Getenv(dockerhostEnvvar)
	if dockerHost == "" {
		return defaultHost
	}

	hostport, err := url.Parse(dockerHost)
	if err != nil {
		return defaultHost
	}

	host, _, err := net.SplitHostPort(hostport.Host)
	if err != nil {
		return defaultHost
	}
	ip := net.ParseIP(host)
	if ip.To4() != nil {
		return host
	}
	//ipv6
	if ip.To16() != nil {
		return fmt.Sprintf("[%v]", host)
	}
	//name
	return host
}
