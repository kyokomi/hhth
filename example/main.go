package main

import "net/http"

func init() {
	http.HandleFunc("/hoge", HogeHandler)
	http.HandleFunc("/hoge.json", HogeJSONHandler)
	http.HandleFunc("/header", HeaderHandler)
}

func main() {
	http.ListenAndServe(":8000", nil)
}

func HogeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("hogehoge"))
}

func HogeJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"name": "hogehoge", "age": 20}`))
}

func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	xAppHoge := r.Header.Get("X-App-Hoge")
	if xAppHoge != "hoge-header" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("header bad reqeust error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("header ok " + xAppHoge))
}
