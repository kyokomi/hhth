package example

import (
	"io/ioutil"
	"net/http"
)

func init() {
	http.HandleFunc("/hoge", hogeHandler)
	http.HandleFunc("/hoge.json", hogeJSONHandler)
}

func hogeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "HEAD":
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hoge"))
	case "PUT":
		message := r.FormValue("message")
		if message == "hello" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(message))
		} else {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"message": "hello", "url": "/hoge?message=hello"}`))
		}
	case "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		w.Header().Set("Location", "http://localhost:8080/hoge/"+string(data))
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	case "DELETE":
		w.WriteHeader(http.StatusNoContent)
	case "OPTIONS":
		w.Header().Add("Allow", "GET,HEAD,PUT,POST,DELETE")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}

func hogeJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"name": "hoge", "age": 20}`))
}
