package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kyokomi/hhth"
)

func TestHogeHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	testCase1 := hhth.NewTestCase(http.StatusOK, "text/plain; charset=utf-8")
	resp := hhtHelper.Get("/hoge").Do(testCase1)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}

func TestHogeJSONHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	testCase := hhth.NewTestCase(http.StatusOK, "application/json; charset=UTF-8")
	var resp map[string]interface{}
	if err := hhtHelper.Get("/hoge.json").Do(testCase).JSON(&resp); err != nil {
		t.Errorf("error %s", err)
	}
	fmt.Println(resp)
}

func TestHogeHeaderHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)
	testCase := hhth.NewTestCase(http.StatusOK, "text/plain; charset=utf-8")
	resp := hhtHelper.SetHeader("X-App-Hoge", "hoge-header").Get("/header").Do(testCase)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	fmt.Println(resp.String())
}
