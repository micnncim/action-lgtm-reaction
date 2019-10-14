package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "INPUT_TRIGGER: %s", os.Getenv("INPUT_TRIGGER"))
}
