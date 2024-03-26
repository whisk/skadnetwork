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
	postback, err := readAllFromFile(filename)
	if err != nil {
		fmt.Println("error while reading postback:", err)
		os.Exit(1)
	}
	validator := skadnetwork.NewPostbackValidator()
	ok, err := validator.Validate(postback)
	if err != nil {
		fmt.Println("error while validating:", err)
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
