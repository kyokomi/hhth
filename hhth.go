package hhth

import (
	"io"
	"net/http"
	"net/http/httptest"
)

const (
	contentTypeJson = "application/json; charset=UTF-8"
)

type HTTPHandlerTestHelper interface {
	SetHeader(key, value string)
	SetForm(key, value string)

	// method
	Get(urlStr string, testCase TestCase) Response
	Post(urlStr string, bodyType string, body io.Reader, testCase TestCase) Response
}

var _ HTTPHandlerTestHelper = (*httpHandlerTestHelper)(nil)

func New(handler http.Handler) HTTPHandlerTestHelper {
	return &httpHandlerTestHelper{
		handler: handler,
		headers: map[string]string{},
		form:    map[string]string{},
	}
}

type httpHandlerTestHelper struct {
	handler http.Handler

	method  string
	url     string
	headers map[string]string
	form    map[string]string
}

func (h *httpHandlerTestHelper) Get(urlStr string, testCase TestCase) Response {
	h.method = "GET"
	h.url = urlStr
	return h.do(testCase, nil)
}

func (h *httpHandlerTestHelper) Post(urlStr string, bodyType string, body io.Reader, testCase TestCase) Response {
	h.method = "POST"
	h.url = urlStr
	h.SetHeader("Content-Type", bodyType)
	return h.do(testCase, body)
}

func (h *httpHandlerTestHelper) SetHeader(key, value string) {
	h.headers[key] = value
}

func (h *httpHandlerTestHelper) SetForm(key, value string) {
	h.form[key] = value
}

func (h *httpHandlerTestHelper) do(testCase TestCase, body io.Reader) *response {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(h.method, h.url, body)
	if body == nil {
		for key, val := range h.form {
			req.Form.Set(key, val)
		}
	}

	for key, val := range h.headers {
		req.Header.Set(key, val)
	}

	if err != nil {
		return &response{err: err, response: nil}
	}

	h.handler.ServeHTTP(resp, req)

	if err := testCase.Execute(resp); err != nil {
		return &response{err: err, response: nil}
	}

	return &response{err: nil, response: resp}
}
