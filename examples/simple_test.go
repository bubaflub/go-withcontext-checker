package examples

import (
	"bytes"
	"net/http"
)

func bad() {
	url := "http://example.com"
	payload := "this should be flagged by the context checker"
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	req.WithContext(nil)
}

func good() {
	url := "http://example.com"
	payload := "this should be flagged by the context checker"
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	req = req.WithContext(nil)
}
