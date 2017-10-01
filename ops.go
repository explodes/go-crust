package crust

type OpCode byte

const (
	OpPutln        = OpCode(1) // (), print '\n' to stdout
	OpDup          = OpCode(2) // (), duplicate the top of the stack
	OpPut          = OpCode(3) // (), consume and print top of stack to stdout
	OpJump         = OpCode(4) // (line:int), jump to line number
	OpJumpLessThan = OpCode(5) // (value:int, line:int), if the consumed top of stack is less than value, jump to line number

	OpIpush     = OpCode(11) // (value:int), push value onto stack
	OpIadd      = OpCode(12) // (), consume top two values of stack, push sum onto stack
	OpIsubtract = OpCode(14) // (), consume top two values of stack, push (top-1) - (top) onto stack

	OpSpush = OpCode(21) // (value:string), push value onto stack
	OpSadd  = OpCode(22) // (), consume top two values of stack, push concatenation onto stack
)

const (
	InstructionPutln        = "putln"
	InstructionDup          = "dup"
	InstructionPut          = "put"
	InstructionJump         = "jump"
	InstructionJumpLessThan = "jumpl"

	InstructionIpush     = "ipush"
	InstructionIadd      = "iadd"
	InstructionIsubtract = "isub"

	InstructionSpush = "spush"
	InstructionSadd  = "sadd"
)

type ArgType int

const (
	argInt    ArgType = iota
	argString
)

type instructionSignature struct {
	op   OpCode
	args []ArgType
}

var (
	instructionSignatures = map[string]instructionSignature{
		InstructionPutln:        {OpPutln, nil},
		InstructionDup:          {OpDup, nil},
		InstructionPut:          {OpPut, nil},
		InstructionJump:         {OpJump, []ArgType{argInt}},
		InstructionJumpLessThan: {OpJumpLessThan, []ArgType{argInt, argInt}},

		InstructionIpush:     {OpIpush, []ArgType{argInt}},
		InstructionIadd:      {OpIadd, nil},
		InstructionIsubtract: {OpIsubtract, nil},

		InstructionSpush: {OpSpush, []ArgType{argString}},
		InstructionSadd:  {OpSadd, nil},
	}
)
