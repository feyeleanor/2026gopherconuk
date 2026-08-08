package main

import "net/http"

const ADDRESS = ":3000"

func main() {
	http.ListenAndServe(
		ADDRESS,
		http.FileServer(
			http.Dir(".")))
}