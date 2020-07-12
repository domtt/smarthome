package main

import (
	"fmt"
	"net/http"
)

func bellHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "written")
}

func main() {
	http.HandleFunc("/api/bell", bellHandler)
	fmt.Println("Server running")
	http.ListenAndServe(":80", nil)
}
