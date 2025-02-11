package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	startStr := "Hello, OTUS!"
	resultStr := reverse.String(startStr)
	fmt.Println(resultStr)
}
