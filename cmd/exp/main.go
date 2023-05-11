package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
}

type UserMeta struct {
	Visits int
}

func main() {
	// Paths are relative to where we run our code from
	t, err := template.ParseFiles("hello_go.html")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Susan Smith",
		Age:  111,
		Meta: UserMeta{
			Visits: 4,
		},
	}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
