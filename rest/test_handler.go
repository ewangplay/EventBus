package rest

import (
	"fmt"
)

type TestHandler struct {
	*BaseHandler
}

func (this *TestHandler) ProcessFunc(method string, resources []string, params map[string]string, body []byte, result map[string]interface{}) error {

	switch method {
	case "POST":
		return this.PaymentNotify(body, result)
	}

	return fmt.Errorf("unsupported http method: %v", method)
}

func (this *TestHandler) PaymentNotify(body []byte, result map[string]interface{}) error {

	this.Debug("[NOTIFIER_TEST] Request body: %v", string(body))

	return nil
}
