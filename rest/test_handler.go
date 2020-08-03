package rest

import (
	"fmt"
)

// TestHandler struct define
type TestHandler struct {
	*BaseHandler
}

// ProcessFunc ...
func (th *TestHandler) ProcessFunc(method string, resources []string, params map[string]string, body []byte, result map[string]interface{}) error {

	switch method {
	case "POST":
		return th.paymentNotify(body, result)
	}

	return fmt.Errorf("unsupported http method: %v", method)
}

func (th *TestHandler) paymentNotify(body []byte, result map[string]interface{}) error {

	th.Debug("[NOTIFIER_TEST] Request body: %v", string(body))

	return nil
}
