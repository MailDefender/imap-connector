package configuration

import "errors"

type AuthenticationConfiguration struct {
	Username string
	Password string
}

type ServerConfiguration struct {
	Host string
	Port int
	Tls  bool
}

type MailServerConfiguration struct {
	Server         ServerConfiguration
	Authentication AuthenticationConfiguration
}

func (s MailServerConfiguration) Check() error {

	if err := s.Server.check(); err != nil {
		return err
	}

	if err := s.Authentication.check(); err != nil {
		return err
	}

	return nil
}

func (s ServerConfiguration) check() error {
	if s.Host == "" {
		return errors.New("host not set")
	}

	if s.Port <= 0 {
		return errors.New("port not set")
	}
	return nil
}

func (s AuthenticationConfiguration) check() error {
	if s.Username == "" {
		return errors.New("username not set")
	}

	if s.Password == "" {
		return errors.New("password not set")
	}
	return nil
}
