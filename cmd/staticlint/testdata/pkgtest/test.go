package pkgtest

import (
	"fmt"
	"os"
)

func test() {
	fmt.Println("os.Exit() but not in package main")
	os.Exit(1)
}
