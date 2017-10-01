package crust

import (
	"io"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"log"
)

// Program is a parsed crust program
type Interpreter struct {
	// program is the parsed crust program
	program *Program

	// ip is the current instruction pointer
	ip int

	// stack is the state of the program
	stack []interface{}

	// stdout is the destination writer for printing information
	stdout io.Writer

	debug bool
}

type InterpreterOption func(*Interpreter)

func NewInterpreter(program *Program, opts ...InterpreterOption) *Interpreter {
	interpreter := &Interpreter{
		program: program,
		ip:      0,
		stack:   make([]interface{}, 0, 64),
		stdout:  os.Stdout,
		debug:   false,
	}
	for _, opt := range opts {
		opt(interpreter)
	}
	return interpreter
}

func EnableDebug(debug bool) InterpreterOption {
	return func(interpreter *Interpreter) {
		interpreter.debug = debug
	}
}

func WithStdout(w io.Writer) InterpreterOption {
	return func(interpreter *Interpreter) {
		interpreter.stdout = w
	}
}

// Run runs the interpreter until completion.
// If an error occurs during execution, that error is returned.
func (i *Interpreter) Run() error {
	for {
		if err := i.Step(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// Step runs the program for a single instruction.
// If there are no more instructions, EOF is returned.
// If an error occurs during execution, that error is returned.
func (i *Interpreter) Step() error {
	instruction, err := i.nextInstruction()
	if err != nil {
		return err
	}
	switch op := instruction.(type) {
	case OpCode:
		err := i.executeOp(op)
		i.dlog("stack: %#v", i.stack)
		return err
	default:
		return errors.Errorf("invalid program, not an op code: %v", instruction)
	}
}

func (i *Interpreter) push(v interface{}) {
	i.stack = append(i.stack, v)
}

func (i *Interpreter) pop() (interface{}, error) {
	if len(i.stack) == 0 {
		return nil, errors.New("stack is empty")
	}
	var value interface{}
	value, i.stack = i.stack[len(i.stack)-1], i.stack[:len(i.stack)-1]
	return value, nil
}

func (i *Interpreter) peek() (interface{}, error) {
	if len(i.stack) == 0 {
		return nil, errors.New("stack is empty")
	}
	return i.stack[len(i.stack)-1], nil
}

func (i *Interpreter) popInt() (int, error) {
	v, err := i.pop()
	if err != nil {
		return 0, err
	}
	return asInt(v)
}

func (i *Interpreter) popString() (string, error) {
	v, err := i.pop()
	if err != nil {
		return "", err
	}
	return asString(v)
}

func (i *Interpreter) executeOp(op OpCode) error {
	switch op {
	case OpPutln:
		i.toStdout("\n")
		i.dlog("putln")
		return nil
	case OpDup:
		value, err := i.peek()
		if err != nil {
			return err
		}
		i.push(value)
		i.dlog("dup %v", value)
		return nil
	case OpPut:
		top, err := i.pop()
		if err != nil {
			return err
		}
		i.toStdout(top)
		i.dlog("put %v", top)
		return nil
	case OpJump:
		line, err := i.nextInt()
		if err != nil {
			return err
		}
		err = i.jump(line)
		i.dlog("jump %d => %d", line, i.ip)
		return err
	case OpJumpLessThan:
		value, err := i.nextInt()
		if err != nil {
			return err
		}
		line, err := i.nextInt()
		if err != nil {
			return err
		}
		top, err := i.popInt()
		if err != nil {
			return err
		}
		if top < value {
			i.jump(line)
		}
		i.dlog("jump %d<%d ? %d => %d jumped=%v", top, value, line, i.ip, top < value)
		return nil
	case OpIpush:
		value, err := i.nextInt()
		if err != nil {
			return err
		}
		i.push(value)
		i.dlog("ipush %d", value)
		return nil
	case OpIadd:
		a, err := i.popInt()
		if err != nil {
			return err
		}
		b, err := i.popInt()
		if err != nil {
			return err
		}
		c := b + a
		i.push(c)
		i.dlog("iadd %d + %d = %d", b, a, c)
		return nil
	case OpIsubtract:
		a, err := i.popInt()
		if err != nil {
			return err
		}
		b, err := i.popInt()
		if err != nil {
			return err
		}
		c := b - a
		i.push(c)
		i.dlog("isub %d - %d = %d", b, a, c)
		return nil
	case OpSpush:
		value, err := i.nextString()
		if err != nil {
			return err
		}
		i.push(value)
		i.dlog("spush %s", value)
		return nil
	case OpSadd:
		a, err := i.popString()
		if err != nil {
			return err
		}
		b, err := i.popString()
		if err != nil {
			return err
		}
		c := b + a
		i.push(c)
		i.dlog("sadd %s + %s = %s", b, a, c)
		return nil
	}
	return errors.Errorf("invalid op code: %v", op)
}

func (i *Interpreter) jump(line int) error {
	index := line - 1 // convert 1-based line number to 0-base jumpTable index
	if index < 0 || index >= len(i.program.jumpTable) {
		return errors.New("invalid jump index")
	}
	i.ip = i.program.jumpTable[index]
	return nil
}

func (i *Interpreter) nextInstruction() (interface{}, error) {
	if i.ip == len(i.program.instructions) {
		return nil, io.EOF
	}
	instruction := i.program.instructions[i.ip]
	i.ip++
	return instruction, nil
}

func (i *Interpreter) nextInt() (int, error) {
	instruction, err := i.nextInstruction()
	if err != nil {
		return 0, err
	}
	return asInt(instruction)
}

func (i *Interpreter) nextString() (string, error) {
	instruction, err := i.nextInstruction()
	if err != nil {
		return "", err
	}
	return asString(instruction)
}

func (i *Interpreter) toStdout(args ...interface{}) (n int, err error) {
	return fmt.Fprint(i.stdout, args...)
}

func (i *Interpreter) toStdoutf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(i.stdout, format, args...)
}

func (i *Interpreter) dlog(format string, args ...interface{}) {
	if !i.debug {
		return
	}
	log.Printf(format, args...)
}

func asInt(v interface{}) (int, error) {
	value, ok := v.(int)
	if !ok {
		return 0, errors.Errorf("value not int: %v", v)
	}
	return value, nil
}

func asString(v interface{}) (string, error) {
	value, ok := v.(string)
	if !ok {
		return "", errors.Errorf("value not string: %v", v)
	}
	return value, nil
}
