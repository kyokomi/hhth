package hhth

import (
	"encoding/json"
	"net/http/httptest"
)

type Response interface {
	Error() error
	String() string
	Result() (*httptest.ResponseRecorder, error)
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
	if r.err != nil {
		return ""
	}
	return r.response.Body.String()
}

func (r *response) JSON(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	if err := json.Unmarshal(r.response.Body.Bytes(), v); err != nil {
		return err
	}
	return nil
}
