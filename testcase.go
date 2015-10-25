package hhth

import "fmt"

type HandlerTestCaseFunc func(resp Response) error

// Execute calls f(resp) error.
func (f HandlerTestCaseFunc) Execute(resp Response) error {
	return f(resp)
}

type HandlerTestCase interface {
	Execute(resp Response) error
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

func (t *contentLengthTestCase) Execute(resp Response) error {
	r, _ := resp.Result()
	if t.contentLength != r.Body.Len() {
		return fmt.Errorf("should Content-Length %d ≠ %d", t.contentLength, r.Body.Len())
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

func (t *statusCodeTestCase) Execute(resp Response) error {
	r, _ := resp.Result()
	if r.Code != t.statusCode {
		return fmt.Errorf("should return HTTP OK %d ≠ %d", r.Code, t.statusCode)
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

func (t *contentTypeTestCase) Execute(resp Response) error {
	r, _ := resp.Result()
	contentType := r.Header().Get("Content-Type")
	if contentType != t.contentType {
		return fmt.Errorf("should Content-Type %s ≠ %s", contentType, t.contentType)
	}
	return nil
}
