package main

import (
	"html/template"
	"os"
)

type User struct {
	Name  string
	Age   int
	Score int
}

func main() {
	// Paths are relative to where we run our code from
	t, err := template.ParseFiles("hello_go.html")
	if err != nil {
		panic(err)
	}

	user := User{
		Name:  "Zozak",
		Age:   111,
		Score: 77,
	}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
