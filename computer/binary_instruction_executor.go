package computer

type riscVBinaryInstructionExecutor struct {
	parser binaryParser
}

type binaryParser interface {
	parse(binary uint32) binaryParseResult
}

type binaryParseResult interface {
	isRType() bool
	isIType() bool
	isSType() bool
	isUType() bool
}

func (ex *riscVBinaryInstructionExecutor) execute(binaryInstruction uint32) {
	parseResult := ex.parser.parse(binaryInstruction)

	if parseResult.isRType() {

	} else if parseResult.isIType() {

	} else if parseResult.isSType() {

	} else if parseResult.isUType() {

	} else {

	}
}
