package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func abort(i interface{}) {
	fmt.Println("Fatal:", i)
	os.Exit(1)
}

var (
	typname = flag.String("type", "", "type to be modified")
)

func main() {
	flag.Parse()
	if *typname == "" {
		abort("no typename provided")
	}
	file := os.Getenv("GOFILE")
	if file == "" {
		abort("no input file specified")
	}
	if err := AddFieldsTo(file, *typname); err != nil {
		abort(err)
	}
	if err := WriteImplFileTo(strings.TrimSuffix(file, ".go") + "_list.go"); err != nil {
		abort(err)
	}
}
