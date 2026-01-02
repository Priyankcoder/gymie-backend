package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "pri@expo"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error generating hash:", err)
	}

	fmt.Println("Bcrypt hash for password 'pri@expo':")
	fmt.Println(string(hash))
}
