package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("os.Exit() in main()")
	os.Exit(1) // want "Detected direct call to os.Exit in the main function of the main package"
}
