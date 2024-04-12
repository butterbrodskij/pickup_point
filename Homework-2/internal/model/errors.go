package model

import "errors"

var (
	ErrorObjectNotFound      = errors.New("object not found")
	ErrorInvalidEnvironment  = errors.New("invalid environment")
	ErrorInvalidInput        = errors.New("invalid input")
	ErrorExcessWeight        = errors.New("excess weight")
	ErrorInvalidKafkaMessage = errors.New("invalid kafka message")
	ErrorEmptyRequest        = errors.New("empty request")
	ErrorEmptyBodyRequest    = errors.New("empty body request")
)
