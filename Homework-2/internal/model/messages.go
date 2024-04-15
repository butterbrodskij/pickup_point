package model

import (
	"net/http"
	"net/url"
	"time"
)

var (
	MessageSuccess = []byte("operation completed successfully")
)

type RequestMessage struct {
	CaughtTime time.Time
	Request    *http.Request
}

type LogMessage struct {
	CaughtTime time.Time `json:"time"`
	Method     string    `json:"method"`
	Url        *url.URL  `json:"url"`
	Body       string    `json:"body"`
	Login      string    `json:"login"`
}
