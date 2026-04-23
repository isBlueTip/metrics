package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
)

func getReader(req *http.Request) (io.Reader, error) {
	if req.Header.Get("Content-Encoding") != "gzip" {
		return req.Body, nil
	}
	reader, err := gzip.NewReader(req.Body)
	if err != nil {
		return nil, err
	}
	return reader, nil
}