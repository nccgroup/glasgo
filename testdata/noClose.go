package main

import(
	"os"
)

func noClose() int {
        file, err := os.Open("noClose.go")
        os.Open("pristAST.go");
	if err != nil {
		os.Exit(1)
	}
	if file != nil {
		return 0
	}
	if _, err := os.Open("printAST.go"); err != nil {
		return 1
	}
	return 0
}
