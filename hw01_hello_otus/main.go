package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const hello = "Hello, OTUS!"

func main() {
	fmt.Println(stringutil.Reverse(hello))
}
