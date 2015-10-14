package testcase

import (
	"fmt"
	"net/http/httptest"
)

type contentTypeTestCase struct {
	contentType string
}

func ContentType(contentType string) HandlerTestCase {
	return &contentTypeTestCase{
		contentType: contentType,
	}
}

var _ HandlerTestCase = (*contentTypeTestCase)(nil)

func (t *contentTypeTestCase) Execute(resp *httptest.ResponseRecorder) error {
	contentType := resp.Header().Get("Content-Type")
	if contentType != t.contentType {
		return fmt.Errorf("should Content-Type %s ≠ %s", contentType, t.contentType)
	}
	return nil
}
