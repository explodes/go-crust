package crust

import (
	"os"
	"github.com/pkg/errors"
	"io"
	"bufio"
	"strconv"
)

// Program is a parsed crust program
type Program struct {
	// instructions is a list of op codes and
	// arguments that describe a program
	instructions []interface{}

	// jumpTable is a mapping between line numbers and
	// their op code position in instructions. The jumpTable is 0-based whereas
	// real line numbers are 1-based
	jumpTable []int
}

// NewProgramFromFile reads a program from disk and creates the program for it
func NewProgramFromFile(path string) (*Program, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open program file")
	}
	defer f.Close()
	return NewProgramFromReader(f)
}

// NewProgramFromReader reads a program from a reader and creates the program for it
func NewProgramFromReader(r io.Reader) (*Program, error) {
	instructions, jumpTable, err := parseProgram(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse r")
	}
	program := &Program{
		instructions: instructions,
		jumpTable:    jumpTable,
	}
	return program, nil
}

func parseProgram(program io.Reader) (instructions []interface{}, jumpTable []int, err error) {
	in := bufio.NewScanner(program)
	in.Split(bufio.ScanWords)

	instructions = make([]interface{}, 0, 64)
	jumpTable = make([]int, 0)
	currentInstructions := new([16]interface{})

	for in.Scan() {
		if err := in.Err(); err != nil {
			return nil, nil, errors.Wrap(err, "unable to scan program")
		}
		token := in.Text()
		n, err := parseOp(token, in, currentInstructions)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to parse op code")
		}
		if n > 0 {
			jumpTable = append(jumpTable, len(instructions))
			instructions = append(instructions, currentInstructions[:n]...)
		}
	}

	return instructions, jumpTable, nil
}

func parseOp(token string, in *bufio.Scanner, instructions *[16]interface{}) (n int, err error) {

	// check for no-argument ops
	signature, ok := instructionSignatures[token]
	if !ok {
		return 0, errors.Errorf("invalid instruction %s", token)
	}

	instructions[0] = signature.op

	for index, argType := range signature.args {
		value, err := getArgument(in, argType)
		if err != nil {
			return 0, err
		}
		instructions[1+index] = value
	}

	n = 1 + len(signature.args)
	return n, nil
}

func getArgument(in *bufio.Scanner, argType ArgType) (interface{}, error) {
	switch argType {
	case argInt:
		return nextInt(in)
	case argString:
		return nextString(in)
	}
	return nil, errors.New("unknown argument type")
}

func nextInt(in *bufio.Scanner) (int, error) {
	if !in.Scan() {
		return 0, errors.New("end of program")
	}
	if err := in.Err(); err != nil {
		return 0, errors.Wrap(err, "unable to advance scanner")
	}
	return strconv.Atoi(in.Text())
}

func nextString(in *bufio.Scanner) (string, error) {
	if !in.Scan() {
		return "", errors.New("end of program")
	}
	if err := in.Err(); err != nil {
		return "", errors.Wrap(err, "unable to advance scanner")
	}
	return in.Text(), nil
}
