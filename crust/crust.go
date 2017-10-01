package main

import (
	"os"
	"github.com/pkg/errors"
	"fmt"
	"github.com/explodes/practice/crustasm"
	"io"
)

func main() {

	if len(os.Args) != 2 {
		exitWith(errors.New("program file not specified"))
	}

	program, err := crustasm.NewProgramFromFile(os.Args[1])
	if err != nil {
		exitWith(errors.Wrap(err, "unable to run program"))
	}

	interpreter := crustasm.NewInterpreter(program, crustasm.EnableDebug(false))
	if err := interpreter.Run(); err != nil {
		if err != io.EOF {
			exitWithCode(2, err)
		}
	}
}

func exitWith(err error) {
	exitWithCode(1, err)
}

func exitWithCode(code int, err error) {
	fmt.Println(err)
	os.Exit(code)
}
