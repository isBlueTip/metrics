package handlers

import (
	"net/http"

	"github.com/isBlueTip/metrics/internal/service"
)

func PingDB(DBConn string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := service.Ping(DBConn); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}
