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

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ Page</h1>
<ul>
  <li>
    <b>Is there a free version?</b>
    Yes! We offer a free trial for 30 days on any paid plans.
  </li>
  <li>
    <b>What are your support hours?</b>
    We have support staff answering emails 24/7, though response
    times may be a bit slower on weekends.
  </li>
  <li>
    <b>How do I contact support?</b>
    Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>
  </li>
</ul>
`)
}

type Router struct {
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "page not found")

		// http.NotFound(w, r)

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func main() {
	/*
		http.Handler => Interface with the ServeHTTP method
		http.HandlerFunc => A function type that has same arguments as ServeHTTP method. Also implements http.Handler interface
	*/

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", Router{})
}
