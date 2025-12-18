package main

import (
	"github.com/xxx/metrics/internal/repository"
	"github.com/xxx/metrics/internal/server"
	"net/http"
)

func run() error {
	mux := server.Router(repository.NewStorage())

	return http.ListenAndServe("0.0.0.0:8080", mux)

}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
