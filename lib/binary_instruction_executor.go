package computer

type riscVBinaryInstructionExecutor struct {
	factory executionFactory
}

type executionFactory interface {
	produce(instruction uint32) executor
}

type executor interface {
	execute()
}

func (ex *riscVBinaryInstructionExecutor) execute(binaryInstruction uint32) {
	(ex.factory.produce(binaryInstruction)).execute()
}
