package computer

type RiscVAssemblyExecutor struct {
	parser assemblyParser
}

type assemblyParser interface {
	parse(assembly string) assemblyParseResult
}

type assemblyParseResult interface {
	isRType() bool
	isIType() bool
	isSType() bool
	isUType() bool
}

func (ex *RiscVAssemblyExecutor) execute(assemblyInstruction string) {
	parseResult := ex.parser.parse(assemblyInstruction)

	if parseResult.isRType() {

	} else if parseResult.isIType() {

	} else if parseResult.isSType() {

	} else if parseResult.isUType() {

	} else {

	}
}
