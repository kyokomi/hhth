package hhth

import (
	"fmt"
	"net/http/httptest"
)

type HandlerTestCase interface {
	Execute(resp *httptest.ResponseRecorder) error
}

type contentLengthTestCase struct {
	contentLength int
}

func TestCaseContentLength(contentLength int) HandlerTestCase {
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

type statusCodeTestCase struct {
	statusCode int
}

func TestCaseStatusCode(statusCode int) HandlerTestCase {
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

type contentTypeTestCase struct {
	contentType string
}

func TestCaseContentType(contentType string) HandlerTestCase {
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
