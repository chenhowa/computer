package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, world!")
	argsWithoutProg := os.Args[1:]
	fmt.Println(len(argsWithoutProg))
	for i := uint(0); i < uint(len(argsWithoutProg)); i++ {
		fmt.Println(argsWithoutProg[i])
	}
}
