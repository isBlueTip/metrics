package main

import (
	"fmt"
	"net"
)

type ServerAddress struct {
	Host string
	Port string
}

func (a *ServerAddress) String() string {
	return fmt.Sprintf("%q:%q\n", a.Host, a.Port)
}

func (a *ServerAddress) Set(s string) error {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}
	a.Port = port
	a.Host = host

	return nil
}
