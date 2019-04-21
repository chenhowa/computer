package computer

type riscVAssemblyParser struct {
}

func (parser *riscVAssemblyParser) parse(instruction string) riscVAssemblyParseResult {

	return riscVAssemblyParseResult{}
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

type riscVAssemblyParseResult struct {
}

/*
	Will need to write a connector between this parser and the assembly executor
	that will use it.
*/
