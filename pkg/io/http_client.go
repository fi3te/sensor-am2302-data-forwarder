package io

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func HttpPostJson(url string, body any) (*http.Response, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", &buf)
	return resp, err
}
