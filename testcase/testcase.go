package testcase

import (
	"net/http/httptest"
)

type HandlerTestCase interface {
	Execute(resp *httptest.ResponseRecorder) error
}
