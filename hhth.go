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

	AddTestCase(testCases ...HandlerTestCase)
	SetTestCase(testCases ...HandlerTestCase)
	AddTestCaseFunc(testCases ...HandlerTestCaseFunc)
	SetTestCaseFunc(testCaseFunc ...HandlerTestCaseFunc)

	// method
	Get(urlStr string, testCases ...HandlerTestCaseFunc) Response
	Post(urlStr string, bodyType string, body io.Reader, testCases ...HandlerTestCaseFunc) Response
}

var _ HTTPHandlerTestHelper = (*httpHandlerTestHelper)(nil)

func New(handler http.Handler) HTTPHandlerTestHelper {
	return &httpHandlerTestHelper{
		handler: handler,

		method:    "",
		url:       "",
		headers:   map[string]string{},
		form:      map[string]string{},
		testCases: []HandlerTestCase{},
	}
}

type httpHandlerTestHelper struct {
	handler http.Handler

	method    string
	url       string
	headers   map[string]string
	form      map[string]string
	testCases []HandlerTestCase
}

func (h *httpHandlerTestHelper) Get(urlStr string, testCases ...HandlerTestCaseFunc) Response {
	h.method = "GET"
	h.url = urlStr
	return h.do(nil, testCases...)
}

func (h *httpHandlerTestHelper) Post(urlStr string, bodyType string, body io.Reader, testCases ...HandlerTestCaseFunc) Response {
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

func (h *httpHandlerTestHelper) SetTestCase(testCases ...HandlerTestCase) {
	h.testCases = testCases
}

func (h *httpHandlerTestHelper) AddTestCase(testCases ...HandlerTestCase) {
	h.testCases = append(h.testCases, testCases...)
}

func (h *httpHandlerTestHelper) SetTestCaseFunc(testCaseFunc ...HandlerTestCaseFunc) {
	dst := make([]HandlerTestCase, len(testCaseFunc))
	for idx, tf := range testCaseFunc {
		dst[idx] = HandlerTestCase(tf)
	}
	h.testCases = dst
}

func (h *httpHandlerTestHelper) AddTestCaseFunc(testCaseFunc ...HandlerTestCaseFunc) {
	for _, tf := range testCaseFunc {
		h.testCases = append(h.testCases, HandlerTestCase(tf))
	}
}

func (h *httpHandlerTestHelper) do(body io.Reader, testCases ...HandlerTestCaseFunc) Response {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(h.method, h.url, body)
	if err != nil {
		return NewErrorResponse(err)
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

	testResp := NewResponse(resp)

	execTestCases := make([]HandlerTestCase, len(h.testCases), len(h.testCases)+len(testCases))
	copy(execTestCases, h.testCases)
	for idx, tc := range testCases {
		execTestCases[idx+len(h.testCases)-1] = HandlerTestCase(tc)
	}

	for _, testCase := range execTestCases {
		if err := testCase.Execute(testResp); err != nil {
			return NewErrorResponse(err)
		}
	}

	return testResp
}
