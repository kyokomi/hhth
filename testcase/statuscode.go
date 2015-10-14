package testcase

import (
	"fmt"
	"net/http/httptest"
)

type statusCodeTestCase struct {
	statusCode int
}

func StatusCode(statusCode int) HandlerTestCase {
	return &statusCodeTestCase{
		statusCode: statusCode,
	}
}

var _ HandlerTestCase = (*statusCodeTestCase)(nil)

func (t *statusCodeTestCase) Execute(resp *httptest.ResponseRecorder) error {
	if resp.Code != t.statusCode {
		return fmt.Errorf("should return HTTP OK %d ≠ %d", resp.Code, t.statusCode)
	}
	return nil
}
