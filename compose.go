package dct

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

//ErrUnknownServiceOrPort is returned if no mapping for a service/port combination was found
var ErrUnknownServiceOrPort = errors.New("service or port not found")

//ErrDockerComposeNotFound is thrown if docker-compose was not found in exec.LookPath
var ErrDockerComposeNotFound = errors.New("docker-compose not found in path")

// haveCompose returns whether the "docker-compose" command was found or not.
func haveCompose() bool {
	_, err := exec.LookPath("docker-compose")
	return err == nil
}

//TestEnv can be used for controlling services through docker-compose
type TestEnv interface {
	StartAll() error
	Service(serviceName string) Service
	StopAll() error
	RemoveAll() error
	ServiceCount() (int, error)
}

//NewComposer initializes a new Composer for automatic tests
func NewComposer(composeFile string) (TestEnv, error) {
	if !haveCompose() {
		return nil, ErrDockerComposeNotFound
	}
	return composer{
		Host:        "127.0.0.1",
		ComposeFile: composeFile,
	}, nil
}

type composer struct {
	Host        string
	ComposeFile string
}

func (c composer) StartAll() error {
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "up", "-d")
	return runCommand(cmd)
}

func (c composer) StopAll() error {
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "stop")
	return runCommand(cmd)
}

func (c composer) RemoveAll() error {
	if err := c.StopAll(); err != nil {
		return err
	}
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "rm", "-f")
	return runCommand(cmd)
}

func runCommand(cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running docker-compose %v\nStdOut: %s\nStdErr: %s\nError: %v\n\n", cmd.Args, stdout.String(), stderr.String(), err)
	}
	return nil
}

func (c composer) Service(serviceName string) Service {
	return service{
		Composer: c,
		Service:  serviceName,
	}
}

func (c composer) ServiceCount() (int, error) {
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "ps", "-q")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return -1, fmt.Errorf("error running docker-compose\nStdOut: %s\nStdErr: %s\nError: %v\n\n", stdout.String(), stderr.String(), err)
	}
	count := bytes.Count(stdout.Bytes(), []byte{'\n'})
	return count, nil

}

func (c *composer) getServicePort(service, port string) (string, error) {
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "port", service, port)
	// fmt.Printf("cmd args:%v\r\n", cmd.Args)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running docker-compose\nStdOut: %s\nStdErr: %s\nError: %v\n\n", stdout.String(), stderr.String(), err)
	}
	portStr := strings.TrimSpace(stdout.String())
	if portStr == "" {
		return "", ErrUnknownServiceOrPort
	}
	portStr = strings.Replace(portStr, "0.0.0.0:", "", 1)

	return portStr, nil
}

func (c *composer) startService(service string) error {
	if err := runCommand(exec.Command("docker-compose", "-f", c.ComposeFile, "create", service)); err != nil {
		return err
	}
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "start", service)
	return runCommand(cmd)
}

func (c *composer) stopService(service string) error {
	var cmd = exec.Command("docker-compose", "-f", c.ComposeFile, "stop", service)
	return runCommand(cmd)
}
