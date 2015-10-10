package hhth

import (
	"encoding/json"
	"fmt"
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
		params: handlerTestParams{
			Headers: map[string]string{},
			Form:    map[string]string{},
		},
	}
}

type httpHandlerTestHelper struct {
	params  handlerTestParams
	handler http.Handler
}

func (h *httpHandlerTestHelper) Get(urlStr string, testCase TestCase) Response {
	h.params.Method = "GET"
	h.params.URL = urlStr
	return h.do(testCase, nil)
}

func (h *httpHandlerTestHelper) Post(urlStr string, bodyType string, body io.Reader, testCase TestCase) Response {
	h.params.Method = "POST"
	h.params.URL = urlStr
	h.SetHeader("Content-Type", bodyType)
	return h.do(testCase, body)
}

func (h *httpHandlerTestHelper) SetHeader(key, value string) {
	h.params.Headers[key] = value
}

func (h *httpHandlerTestHelper) SetForm(key, value string) {
	h.params.Form[key] = value
}

type handlerTestParams struct {
	Method  string
	URL     string
	Headers map[string]string
	Form    map[string]string
}

type Response interface {
	Error() error
	String() string
	JSON(v interface{}) error
}

type response struct {
	err      error
	response *httptest.ResponseRecorder
}

func (r *response) Error() error {
	return r.err
}

func (r *response) Result() (*httptest.ResponseRecorder, error) {
	return r.response, r.err
}

func (r *response) String() string {
	if r.response == nil {
		return ""
	}
	return r.response.Body.String()
}

func (r *response) JSON(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.response == nil {
		return fmt.Errorf("response is nil")
	}

	if err := json.Unmarshal(r.response.Body.Bytes(), v); err != nil {
		return err
	}
	return nil
}

func (h *httpHandlerTestHelper) do(testCase TestCase, body io.Reader) *response {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(h.params.Method, h.params.URL, body)
	if body == nil {
		for key, val := range h.params.Form {
			req.Form.Set(key, val)
		}
	}

	for key, val := range h.params.Headers {
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
