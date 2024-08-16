package hue

import "net/http"

type RequestType int

const (
	Dimming RequestType = 0
	On      RequestType = 1
	Fetch   RequestType = 2
)

type Request struct {
	requestType RequestType
	request     http.Request
}
