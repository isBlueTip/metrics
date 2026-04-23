package main

import (
	"fmt"
	"net"
	"time"
)

type Config struct {
	Address         *string `env:"ADDRESS"`
	StoreInterval   *int    `env:"STORE_INTERVAL"`
	FileStorePath   *string `env:"FILE_STORAGE_PATH"`
	Restore        *bool   `env:"RESTORE"`
}

type StoreConfig struct {
	Interval     time.Duration
	FilePath     string
	Restore      bool
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
