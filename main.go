package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	context := NewContext()
	parser := flags.NewParser(&context.Option, flags.Default)
	parser.Name = "git-user"
	_, e := parser.Parse()
	if e != nil {
		return
	}
	if parser.Active == nil {
		parser.WriteHelp(os.Stdout)
		return
	}
	if err := context.Execute(parser.Active.Name); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
