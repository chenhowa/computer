package computer

import "fmt"

/*
This file is tightly coupled with the riscVBinaryInstructionParser
*/

type riscVBinaryInstructionExecutionFactory struct {
}

type binaryExecutor interface {
	execute()
}

func (factory *riscVBinaryInstructionExecutionFactory) produce(instruction uint32) binaryExecutor {
	parser := riscVBinaryInstructionParser{}
	result := parser.parse(instruction)

	var executor binaryExecutor

	switch result.instructionType {
	case I:
		executor = factory.produceI(instruction)
	case R:
		executor = factory.produceR(instruction)
	case J:
		executor = factory.produceJ(instruction)
	case B:
		executor = factory.produceB(instruction)
	case S:
		executor = factory.produceS(instruction)
	default:
		panic(fmt.Sprintf("unrecognized instruction type: %d", result.instructionType))
	}

	return executor
}

type executorI struct {
}

func (ex *executorI) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceI(instruction uint32) binaryExecutor {

	var executor = executorI{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorR struct {
}

func (ex *executorR) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceR(instruction uint32) binaryExecutor {

	var executor = executorR{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorJ struct {
}

func (ex *executorJ) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceJ(instruction uint32) binaryExecutor {

	var executor = executorJ{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorB struct {
}

func (ex *executorB) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceB(instruction uint32) binaryExecutor {

	var executor = executorB{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorS struct {
}

func (ex *executorS) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceS(instruction uint32) binaryExecutor {

	var executor = executorS{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}
