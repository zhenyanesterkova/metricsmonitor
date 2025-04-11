package main

import (
	"fmt"
	"os"
)

func NotMain() {
	fmt.Println("os.Exit() but not in main()")
	os.Exit(1)
}
