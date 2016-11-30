package main

import (
	"fmt"

	"./v1"
	"./v2"
)

func main() {
	fmt.Println(v1.Version())
	fmt.Println(v2.Version())
}
