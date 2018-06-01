package main

import (
	"os"
	// "fmt"
	// "io/ioutil"
)

func main() {
	args := os.Args
	ParseFile(args[1])
}
