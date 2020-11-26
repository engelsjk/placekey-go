package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println(os.Getenv("PLACEKEY_API_KEY"))
}
