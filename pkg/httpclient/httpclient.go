package httpclient

import (
	"net/http"
	"time"
)

func New() *http.Client {
	return &http.Client{
		Timeout: 500 * time.Millisecond,
	}
}
