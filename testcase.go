package hhth

import (
	"fmt"
	"net/http/httptest"
)

type TestCase struct {
	statusCode         int    // Required
	contentType        string // Required
	validContentLength bool
	contentLength      int
}

func NewTestCase(statusCode int, contentType string) TestCase {
	return TestCase{
		statusCode:         statusCode,
		contentType:        contentType,
		validContentLength: false,
		contentLength:      0,
	}
}

func (t *TestCase) SetContentLength(contentLength int) {
	t.validContentLength = true
	t.contentLength = contentLength
}

func (t *TestCase) Execute(resp *httptest.ResponseRecorder) error {
	if resp.Code != t.statusCode {
		return fmt.Errorf("should return HTTP OK %d ≠ %d", resp.Code, t.statusCode)
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != t.contentType {
		return fmt.Errorf("should Content-Type %s ≠ %s", contentType, t.contentType)
	}

	if t.validContentLength {
		if t.contentLength != resp.Body.Len() {
			return fmt.Errorf("should Content-Length %s ≠ %s", t.contentLength, resp.Body.Len())
		}
	}

	return nil
}
