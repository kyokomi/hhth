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

	testCase1 := hhth.NewTestCase(http.StatusOK, "text/plain; charset=utf-8")
	resp := hhtHelper.Get("/hoge", testCase1)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestHogeJSONHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	testCase := hhth.NewTestCase(http.StatusOK, "application/json; charset=UTF-8")
	var resp map[string]interface{}
	if err := hhtHelper.Get("/hoge.json", testCase).JSON(&resp); err != nil {
		t.Errorf("error %s", err)
	}
	fmt.Println(resp)
}

func TestHogeHeaderHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	hhtHelper.SetHeader("X-App-Hoge", "hoge-header")
	testCase := hhth.NewTestCase(http.StatusOK, "text/plain; charset=utf-8")
	resp := hhtHelper.Get("/header", testCase)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestPostHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	testCase := hhth.NewTestCase(http.StatusOK, "text/plain; charset=utf-8")

	formData := url.Values{}
	formData.Set("name", "hoge")
	formData.Set("age", "19")

	body := bytes.NewBufferString(formData.Encode())
	resp := hhtHelper.Post("/post", "application/x-www-form-urlencoded", body, testCase)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func init() {
	http.HandleFunc("/hoge", hogeHandler)
	http.HandleFunc("/hoge.json", hogeJSONHandler)
	http.HandleFunc("/header", headerHandler)
	http.HandleFunc("/post", postHandler)
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

func hogeJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderError(http.StatusMethodNotAllowed, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"name": "hogehoge", "age": 20}`))
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
	w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
}
