package main

import (
	"fmt"
	"net"
)

type Config struct {
	Address *string `env:"ADDRESS"`
}

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
	a.Host = host
	a.Port = port

	return nil
}

//func parseFlags() {
//	flag.StringVar()
//}
