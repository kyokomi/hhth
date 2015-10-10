package hhth

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
)

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
