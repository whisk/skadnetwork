package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	skadnetwork "github.com/whisk/skadnetwork/pkg"
)

func main() {
	pflag.Parse()
	if pflag.NArg() < 1 {
		fmt.Printf("Usage:\n")
		fmt.Printf("\t%s <FILE>\n", filepath.Base(os.Args[0]))
		pflag.PrintDefaults()
		os.Exit(1)
	}
	filename := pflag.Arg(0)
	jsonBytes, err := readAllFromFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading postback: %s\n", err)
		os.Exit(1)
	}
	validator := skadnetwork.NewPostbackValidator()
	ok, err := validator.Validate(jsonBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error validating postback: %s\n", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("Postback is not valid. Errors found:")
		for _, e := range validator.Errors() {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}
	fmt.Println("Postback is valid")
}

func readAllFromFile(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
