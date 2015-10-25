package example_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kyokomi/hhth"
)

func TestHogeHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	hhtHelper.SetTestCase(
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
	)

	resp := hhtHelper.Get("/hoge")
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	if resp.String() != "hoge" {
		t.Errorf("error response body hoge != %s", resp.String())
	}
}

func TestJSONParse(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	hhtHelper.SetTestCase(
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("application/json; charset=UTF-8"),
	)

	var resp map[string]interface{}
	if err := hhtHelper.Get("/hoge.json").JSON(&resp); err != nil {
		t.Errorf("error %s", err)
	}

	if resp["name"].(string) != "hoge" {
		t.Errorf("error json response name != %s", resp["name"])
	}

	if resp["age"].(float64) != 20 {
		t.Errorf("error json response age != %s", resp["age"])
	}
}

func TestOptionsHogeHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	hhtHelper.SetTestCase(
		hhth.TestCaseStatusCode(http.StatusNoContent),
		hhth.HandlerTestCaseFunc(func(resp hhth.Response) error {
			r, _ := resp.Result()
			if r.Header().Get("Allow") != "GET,HEAD,PUT,POST,DELETE" {
				return fmt.Errorf("allow header error %s", r.Header().Get("Allow"))
			}
			return nil
		}),
	)

	resp := hhtHelper.Options("/hoge")
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}
