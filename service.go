package dct

import (
	"fmt"
)

//Service is used for controlling a single service
type Service interface {
	Port(port string) (string, error)
	URI(port string) (string, error)
	Start() error
	Stop() error
}

type service struct {
	Composer composer
	Service  string
}

func (s service) Port(port string) (string, error) {
	return s.Composer.getServicePort(s.Service, port)
}

func (s service) URI(port string) (string, error) {
	p, err := s.Composer.getServicePort(s.Service, port)
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf("%v:%v", s.Composer.Host, p)
	return uri, nil
}

func (s service) Start() error {
	return s.Composer.startService(s.Service)
}
func (s service) Stop() error {
	return s.Composer.stopService(s.Service)
}
