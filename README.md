hhth
=====================

[![Circle CI](https://circleci.com/gh/kyokomi/hhth.svg?style=svg)](https://circleci.com/gh/kyokomi/hhth)
[![Coverage Status](https://coveralls.io/repos/kyokomi/hhth/badge.svg?branch=master&service=github)](https://coveralls.io/github/kyokomi/hhth?branch=master)


hhth is httpHandler test helper library of the golang.

## Install

```
go get github.com/kyokomi/hhth
```

## Example

### Handler

```go
package example

import "net/http"

func init() {
	http.HandleFunc("/hoge", hogeHandler)
	http.HandleFunc("/hoge.json", hogeJSONHandler)
}

func hogeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("hoge"))
}

func hogeJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"name": "hoge", "age": 20}`))
}
```

### Get Test

```go
package example_test

import (
	"net/http"
	"testing"

	"github.com/kyokomi/hhth"
)

func TestHogeHandler(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	resp := hhtHelper.Get("/hoge",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("text/plain; charset=utf-8"),
	)
	if resp.Error() != nil {
		t.Errorf("error %s", resp.Error())
	}
	if resp.String() != "hoge" {
		t.Errorf("error response body hoge != %s", resp.String())
	}
}
```

### JSON Parse

```go
package example_test

import (
	"net/http"
	"testing"

	"github.com/kyokomi/hhth"
)

func TestJSONParse(t *testing.T) {
	hhtHelper := hhth.New(http.DefaultServeMux)

	var resp map[string]interface{}
	if err := hhtHelper.Get("/hoge.json",
		hhth.TestCaseStatusCode(http.StatusOK),
		hhth.TestCaseContentType("application/json; charset=UTF-8"),
	).JSON(&resp); err != nil {
		t.Errorf("error %s", err)
	}

	if resp["name"].(string) != "hoge" {
		t.Errorf("error json response name != %s", resp["name"])
	}

	if resp["age"].(float64) != 20 {
		t.Errorf("error json response age != %s", resp["age"])
	}
}
```

## TODO

- [x] test
- [x] example
- [x] POST
- [ ] PUT
- [ ] DELETE
- [ ] OPTIONS


