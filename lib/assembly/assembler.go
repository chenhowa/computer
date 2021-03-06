package assembly

/*RiscVAssembler works to assemble one or more lines of a RiscV assembly program into the
corresponding binary instructions. It never assumes that the entire assembly program has been
given to it; thus, it's internal state depends heavily on the order in which assembly instructions are passed
to the assembler.

For example, internally the assembler maintains a 'symbol table' of labels that have been passed to it.
Passing the assembler an instruction that references a label that hasn't yet been encountered will generate
an error, and the instruction will be rejected (the internal state of the assembler will not change)
*/
type RiscVAssembler struct {
	lineCount LineCount
	tokenizer tokenizer
	parser    parser
	gen       codeGenerator
}

/*Assemble takes a set of string RISC-V `instructions` and converts them into
32-bit machine code instructions. Errors may occur during assembly*/
func (assembler *RiscVAssembler) Assemble(instructions string) ([]uint32, error) {
	tokenStream, errTokens := assembler.tokenizer.Tokenize(instructions)
	if errTokens != nil {
		return nil, errTokens
	}

	tree, count, errParse := assembler.parser.Parse(tokenStream)
	if errParse != nil {
		return nil, errParse
	}
	assembler.lineCount += count // THIS MIGHT NOT BE CORRECT (off by one)

	binInstructions, errInstructions := assembler.gen.Generate(tree)
	if errInstructions != nil {
		return nil, errInstructions
	}

	return binInstructions, nil
}

/*LineCount is simply how many lines have been encountered so far*/
type LineCount uint

/*CharCount is simply a count of chars that have been encountered so far*/
type CharCount uint

type token interface {
	GetTokenType() TokenType
	GetTokenString() string
	GetCharCountSinceNewline() CharCount
}

type tokenStream interface {
	HasNext() bool
	Next() (token, error)
	Save() tokenStreamReset
}

type tokenStreamReset interface {
	Reset()
}

type tokenizer interface {
	Tokenize(tokens string) (tokenStream, error)
}

type parser interface {
	Parse(tokenStream tokenStream) (tree abstractSyntaxTree, linesEncountered LineCount, err error)
}

type abstractSyntaxTree interface {
	GetRootIterator() astIterator
}

type astIterator interface {
	GetNumChildren() uint
	GetAstNode() astNode
	GetParentIterator() (astIterator, error)
	GetChildIterator(index uint) (astIterator, error)
}

type astNode interface {
	GetLineCount()
	GetCharCountSinceNewline() CharCount
	GetTokenType() TokenType
	GetTokenString() string
}

type codeGenerator interface {
	Generate(tree abstractSyntaxTree) ([]uint32, error)
}

/*AstNodeKind is an alias for uint, representing the kinds of nodes*/
type AstNodeKind uint

/*These constants represent the kinds of AstNodes that exist*/
const (
	Token AstNodeKind = iota
	Instruction
	Instructions
	Program
)
