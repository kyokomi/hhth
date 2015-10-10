package hhth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	contentTypeJson = "application/json; charset=UTF-8"
)

type HTTPHandlerTestHelper interface {
	SetHeader(key, value string) HTTPHandlerTestHelper
	SetForm(key, value string) HTTPHandlerTestHelper

	// method
	Get(urlStr string) HTTPHandlerTestDoHelper
}

var _ HTTPHandlerTestHelper = (*httpHandlerTestHelper)(nil)

type HTTPHandlerTestDoHelper interface {
	Do(testCase TestCase) Response
}

var _ HTTPHandlerTestDoHelper = (*httpHandlerTestHelper)(nil)

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

func (h *httpHandlerTestHelper) Do(testCase TestCase) Response {
	return h.do(testCase)
}

func (h *httpHandlerTestHelper) Get(urlStr string) HTTPHandlerTestDoHelper {
	h.params.Method = "GET"
	h.params.URL = urlStr
	return h
}

func (h *httpHandlerTestHelper) SetHeader(key, value string) HTTPHandlerTestHelper {
	h.params.Headers[key] = value
	return h
}

func (h *httpHandlerTestHelper) SetForm(key, value string) HTTPHandlerTestHelper {
	h.params.Form[key] = value
	return h
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

func (h *httpHandlerTestHelper) do(testCase TestCase) *response {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(h.params.Method, h.params.URL, nil)

	values := url.Values{}
	for key, val := range h.params.Form {
		values.Add(key, val)
	}
	req.Form = values

	for key, val := range h.params.Headers {
		req.Header.Add(key, val)
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
