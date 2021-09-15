package main

import "fmt"
import (
	"golang.org/x/example/stringutil"
)

func main() {
	text := "Hello, OTUS!"
	textReversed := stringutil.Reverse(text)

	fmt.Println(textReversed)
}
