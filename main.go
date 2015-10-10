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
	Do(params HandlerTestParams, testCase TestCase, v interface{}) error
}

type httpHandlerTestHelper struct {
	Handler http.Handler
}

func NewHTTPHandlerTestHelper(handler http.Handler) HTTPHandlerTestHelper {
	return &httpHandlerTestHelper{
		Handler: handler,
	}
}

type HandlerTestParams struct {
	Method  string // Required
	URL     string // Required
	Headers map[string]string
	Params  map[string]string
}

type TestCase struct {
	StatusCode  int    // Required
	ContentType string // Required
}

// TODO: not supported POST,DELETE,OPTIONS,PUT
func (h *httpHandlerTestHelper) Do(params HandlerTestParams, testCase TestCase, v interface{}) error {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(params.Method, params.URL, nil)

	values := url.Values{}
	for key, val := range params.Params {
		values.Add(key, val)
	}
	req.Form = values

	for key, val := range params.Headers {
		req.Header.Add(key, val)
	}

	if err != nil {
		return err
	}

	h.Handler.ServeHTTP(resp, req)

	// test

	if resp.Code != testCase.StatusCode {
		return fmt.Errorf("should return HTTP OK %d ≠ %d", resp.Code, testCase.StatusCode)
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != testCase.ContentType {
		return fmt.Errorf("should Content-Type %s ≠ %s", contentType, testCase.ContentType)
	}

	// parse

	// TODO: json parse test
	if contentType == contentTypeJson {
		err := json.Unmarshal(resp.Body.Bytes(), v)
		if err != nil {
			return err
		}
	} else {
		v = resp.Body.String()
	}

	return nil
}
