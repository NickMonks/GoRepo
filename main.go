package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
