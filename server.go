package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
