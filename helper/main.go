package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		os.Exit(0)
	}

	nums := os.Args[1]
	split := strings.Split(nums, " ")
	for _, x := range split {
		fmt.Print("\"", x, "\", ")
	}
}
