package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Println("Invalid command", os.Args[1])
	}
}

func hash(password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing: ", password)
		return
	}

	fmt.Println(string(hashedPassword))
}

func compare(password, hash string) {
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		fmt.Println("Password is invalid: ", password)
		return
	}

	fmt.Println("Password is correct!")
}
