package hhth

import (
	"io"
	"net/http"
	"net/http/httptest"

	"net/url"
)

const (
	contentTypeJson = "application/json; charset=UTF-8"
)

type HTTPHandlerTestHelper interface {
	SetHeader(key, value string)
	SetForm(key, value string)

	// method
	Get(urlStr string, testCases ...HandlerTestCase) Response
	Post(urlStr string, bodyType string, body io.Reader, testCases ...HandlerTestCase) Response
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

func (h *httpHandlerTestHelper) Get(urlStr string, testCases ...HandlerTestCase) Response {
	h.method = "GET"
	h.url = urlStr
	return h.do(nil, testCases...)
}

func (h *httpHandlerTestHelper) Post(urlStr string, bodyType string, body io.Reader, testCases ...HandlerTestCase) Response {
	h.method = "POST"
	h.url = urlStr
	h.SetHeader("Content-Type", bodyType)
	return h.do(body, testCases...)
}

func (h *httpHandlerTestHelper) SetHeader(key, value string) {
	h.headers[key] = value
}

func (h *httpHandlerTestHelper) SetForm(key, value string) {
	h.form[key] = value
}

func (h *httpHandlerTestHelper) do(body io.Reader, testCases ...HandlerTestCase) *response {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(h.method, h.url, body)
	if err != nil {
		return &response{err: err, response: nil}
	}

	if req.Form == nil {
		req.Form = make(url.Values)
	}

	for key, val := range h.form {
		req.Form.Set(key, val)
	}

	for key, val := range h.headers {
		req.Header.Set(key, val)
	}

	h.handler.ServeHTTP(resp, req)

	for _, testCase := range testCases {
		if err := testCase.Execute(resp); err != nil {
			return &response{err: err, response: nil}
		}
	}

	return &response{err: nil, response: resp}
}
