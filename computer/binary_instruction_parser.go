package computer

type riscVBinaryInstructionParser struct {
}

func (parser *riscVBinaryInstructionParser) parse(instruction uint32) riscVBinaryParseResult {
	opcode := 127 & instruction
	var instructionType instructionType

	if opcode < 3 { // immediate arithmetic
		instructionType = i
	} else if opcode < 4 { // register-register arithmetic
		instructionType = r
	} else if opcode < 6 { // jump
		instructionType = j
	} else if opcode < 7 { // branch
		instructionType = b
	} else if opcode < 8 { // loads
		instructionType = i
	} else if opcode < 9 { // stores
		instructionType = s
	} else if opcode < 10 { // system (csr)
		instructionType = i
	}

	return riscVBinaryParseResult{}
}

type instructionType uint

const (
	// the below instruction types are the base instruction types
	r instructionType = iota + 1
	i
	s
	u

	// the below constants are extension instruction types. see the online specification for more.
	m
	a
	f
	d
	q
	l
	c
	b
	j
	t
	p
	v
	n
)

type riscVBinaryParseResult struct {
}

/*
	Will need to write a connector between this parser and the assembly executor
	that will use it.
*/
