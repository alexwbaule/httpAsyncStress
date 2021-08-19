package structs

import (
	"net/http"
	"time"
)

type RequestData struct {
    Query  []struct {
		Name string `json:"name"`
		Values []string `json:"values"`
	}`json:"query"`
    Headers  []struct {
		Name string `json:"name"`
		Value string `json:"value"`
	}`json:"headers"`
	Replaces []struct {
		Values []int  `json:"values,omitempty"`
		Value  int  `json:"value,omitempty"`
		Type   string `json:"type"`
		Sort   string `json:"sort,omitempty"`
		Mark   string `json:"mark"`
		Format string `json:"format,omitempty"`
	}`json:"replaces"`
    Body string   `json:"body"`
	Url  string `json:"url"`
	Type string `json:"type"`
	Grep string `json:"grep"`
}

type LogData struct{
	Key string
	Value string
}
type SoapData struct {
	Rpl []LogData
}

type HttpResult struct {
	Response *http.Response
	Err      error
	Start	time.Time
	Tget	time.Duration
	Worker  int
	Count	int
	Soap	SoapData
}

type HttpRequest struct {
	HttpData RequestData
	Count int
	HttpTimeout time.Duration
}
