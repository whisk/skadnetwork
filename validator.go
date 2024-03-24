package main

import (
	"fmt"
	"io"
	"os"

	skadnetwork "github.com/whisk/skadnetwork/pkg"
)

func main() {
	validator := skadnetwork.NewPostbackValidator()
	postback, err := readAllFromFile("testdata/com.example-4.0.json")
	if err != nil {
		fmt.Println("error while reading postback: ", err)
		os.Exit(1)
	}
	ok, err := validator.Validate(postback)
	if err != nil {
		fmt.Println("error while validating: ", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("postback is not valid. Errors found:")
		for i, e := range validator.Errors() {
			fmt.Printf("#%d. %s\n", i+1, e.String())
		}
		os.Exit(1)
	}
	fmt.Println("postback is valid")
}

func readAllFromFile(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
