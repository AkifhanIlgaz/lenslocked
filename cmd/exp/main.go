package main

import (
	"html/template"
	"os"
)

type User struct {
	Name    string
	Age     int
	Scores  []float64
	Aliases map[string]string
}

func main() {
	// Paths are relative to where we run our code from
	t, err := template.ParseFiles("hello_go.html")
	if err != nil {
		panic(err)
	}

	user := User{
		Name:   "Zozak",
		Age:    111,
		Scores: []float64{1.2, 3.4, 0.5, 77.3},
		Aliases: map[string]string{
			"Zozak":    "ZZK",
			"Starknet": "STRK",
			"Aptos":    "APT",
		},
	}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
