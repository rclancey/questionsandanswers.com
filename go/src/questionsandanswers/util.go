package main

import (
	"io"
	"net/http"
	"strconv"
)

func readAtMost(r *http.Request, size int) ([]byte, error) {
	contentLength, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil || contentLength <= 0 || contentLength > size {
		contentLength = size
	}
	data := make([]byte, contentLength)
	offset := 0
	for {
		n, err := r.Body.Read(data[offset:])
		if n == 0 || err == io.EOF {
			offset += n
			return data[:offset], nil
		}
		if err != nil {
			return nil, err
		}
		offset += n
	}
	return data[:offset], nil
}
