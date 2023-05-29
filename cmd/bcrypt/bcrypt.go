package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		fmt.Println(hash(os.Args[2]))
	case "compare":
		fmt.Println(compare(os.Args[2], os.Args[3]))
	default:
		fmt.Println("Invalid command", os.Args[1])
	}
}

func hash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing: ", password)
		return ""
	}

	return string(hashedPassword)
}

func compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
