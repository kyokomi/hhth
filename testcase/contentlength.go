package testcase

import (
	"fmt"
	"net/http/httptest"
)

type contentLengthTestCase struct {
	contentLength int
}

func ContentLength(contentLength int) HandlerTestCase {
	return &contentLengthTestCase{
		contentLength: contentLength,
	}
}

var _ HandlerTestCase = (*contentLengthTestCase)(nil)

func (t *contentLengthTestCase) Execute(resp *httptest.ResponseRecorder) error {
	if t.contentLength != resp.Body.Len() {
		return fmt.Errorf("should Content-Length %d ≠ %d", t.contentLength, resp.Body.Len())
	}
	return nil
}
