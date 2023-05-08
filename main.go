package main

import (
	"fmt"
	"net/http"
)

// Check https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers for all available HTTP headers

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Header is just a map.  type Header map[string][]string
	// Set() replaces any existing values for the key, Add() appends to existing values
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<h1>Welcome to my awesome site</h1>`)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
	 <h1>
	 Contact Page 
	  </h1>

	 <p>To get in touch, email me at <a href="mailto:akifhanilgaz@gmail.com">akifhanilgaz@gmail.com</a></p>
	
	`)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/contact", contactHandler)
	fmt.Println("Starting the server on :3000")

	http.ListenAndServe(":3000", mux)
}
