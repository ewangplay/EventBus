package i

import (
	"net/http"
)

// ProcessFunc prototype define
type ProcessFunc func(
	method string,
	resources []string,
	params map[string]string,
	body []byte,
	result map[string]interface{}) error

// Handler interface define
type Handler interface {
	Process(
		r *http.Request,
		resources []string,
		f ProcessFunc) (string, error)

	ProcessFunc(
		method string,
		resources []string,
		params map[string]string,
		body []byte,
		result map[string]interface{}) error
}
