package model

import (
	"net/http"
	"time"
)

var (
	MessageSuccess = []byte("operation completed successfully")
)

type RequestMessage struct {
	CaughtTime time.Time
	Request    *http.Request
}
