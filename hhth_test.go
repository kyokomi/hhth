package hhth_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/kyokomi/hhth"
)

func TestHogeHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	resp := hhtHelper.Get("/hoge",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
		hhth.TestCaseContentLength(len("hogehoge")),
	)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestError(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	respError1 := hhtHelper.Get("/hoge",
		hhth.TestCaseStatusCode(http.StatusFound), // error
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
		hhth.TestCaseContentLength(len("hogehoge")),
	)
	if respError1.Error() == nil {
		t.Error("error not error")
	}

	respError2 := hhtHelper.Get("/hoge",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("application/json; charset=UTF-8"), // error
		hhth.TestCaseContentLength(len("hogehoge")),
	)
	if respError2.Error() == nil {
		t.Error("error not error")
	}

	respError3 := hhtHelper.Get("/hoge",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
		hhth.TestCaseContentLength(len("hogehogeaaaaa")), // error
	)
	if respError3.Error() == nil {
		t.Error("error not error")
	}
}

func TestHogeJSONHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	resp := hhtHelper.Get("/hoge.json",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("application/json; charset=UTF-8"),
	)

	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}

	var respJson map[string]interface{}
	if err := resp.JSON(&respJson); err != nil {
		t.Errorf("error %s", err)
	}

	fmt.Println(resp)
}

func TestHogeHeaderHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	hhtHelper.SetHeader("X-App-Hoge", "hoge-header")

	resp := hhtHelper.Get("/header",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
	)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestPostHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	formData := url.Values{}
	formData.Set("name", "hoge")
	formData.Set("age", "19")

	resp := hhtHelper.Post("/post", "application/x-www-form-urlencoded",
		bytes.NewBufferString(formData.Encode()),
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
	)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestGetFormHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	hhtHelper.SetForm("name", "hoge")

	resp := hhtHelper.Get("/get-form",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
	)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestNewRequestError(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	resp := hhtHelper.Get("%||~=&&%%%&&T''=0I)=)((&'") // error
	if resp.Error() == nil {
		t.Error("error not nil")
	}

	var parseJSON map[string]string
	if err := resp.JSON(&parseJSON); err == nil { // error
		t.Error("error not nil")
	}
	if resp.String() != "" { // error
		t.Error("error not nil")
	}

	if res, err := resp.Result(); res != nil || err == nil {
		t.Error("error not nil")
	}
}

func TestJSONParseError(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	resp := hhtHelper.Get("error.json")
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}

	var parseJSON map[string]string
	if err := resp.JSON(&parseJSON); err == nil { // error
		t.Error("error not nil")
	}
}

func init() {
	http.HandleFunc("/hoge", hogeHandler)
	http.HandleFunc("/hoge.json", hogeJSONHandler)
	http.HandleFunc("/error.json", errorJSONHandler)
	http.HandleFunc("/header", headerHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get-form", getFormHandler)
}

func hogeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("hogehoge"))
}

func getFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	if r.FormValue("name") != "hoge" {
		renderError(http.StatusBadRequest, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("get-form"))
}

func hogeJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"name": "hogehoge", "age": 20}`))
}

func errorJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"name": "hogehoge", "age": 20`)) // json parse error
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	xAppHoge := r.Header.Get("X-App-Hoge")
	if xAppHoge != "hoge-header" {
		renderError(http.StatusBadRequest, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("header ok " + xAppHoge))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderError(http.StatusBadRequest, w)
		return
	}

	if r.PostForm.Encode() != "age=19&name=hoge" {
		renderError(http.StatusBadRequest, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("post ok " + r.PostForm.Encode()))
}

func renderError(statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}
