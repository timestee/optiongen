package main

import (
	"fmt"
	"os"

	"github.com/sandwich-go/xconf"
	"github.com/timestee/optiongen"
)

func main() {
	xconf.Parse(optiongen.AtomicConfig(), xconf.WithLogDebug(func(string) {}), xconf.WithEnviron())
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("unable to get working directory: %v", err))
	}
	optiongen.ParseDir(wd)
}
