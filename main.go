package main

import (
	"fmt"
	"github.com/DennisDenuto/uber-client/cli/cmd"
	goflags "github.com/jessevdk/go-flags"
	"os"
)

func main() {
	opts := &cmd.Opts{}
	parser := goflags.NewParser(opts, goflags.HelpFlag|goflags.PassDoubleDash)
	_, err := parser.ParseArgs(os.Args[1:])
	_, ok := err.(*goflags.Error)

	if ok {
		fmt.Println(err)
	}
}
