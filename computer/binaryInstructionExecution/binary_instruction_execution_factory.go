package binaryInstructionExecution

import "fmt"

/*
This file is tightly coupled with the riscVBinaryInstructionParser
*/

type riscVBinaryInstructionExecutionFactory struct {
	executor riscVExecutor
}

type riscVExecutor interface {
	addImmediate(dest uint, reg uint, immediate uint32)
	setLessThanImmediate(dest uint, reg uint, immediate uint32)
	setLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32)
	andImmmediate(dest uint, reg uint, immediate uint32)
	orImmediate(dest uint, reg uint, immediate uint32)
	xorImmediate(dest uint, reg uint, immediate uint32)
	shiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32)
	shiftRightLogicalImmediate(dest uint, reg uint, immediate uint32)
	shiftRightArithmeticImmediate(dest uint, reg uint, immediate uint32)
	loadUpperImmediate(dest uint, immediate uint32)
	addUpperImmediateToPC(dest uint, immediate uint32)
	add(dest uint, reg1 uint, reg2 uint)
	setLessThan(dest uint, reg1 uint, reg2 uint)
	setLessThanUnsigned(dest uint, reg1 uint, reg2 uint)
	and(dest uint, reg1 uint, reg2 uint)
	or(dest uint, reg1 uint, reg2 uint)
	xor(dest uint, reg1 uint, reg2 uint)
	shiftLeftLogical(dest uint, reg uint, shiftreg uint)
	shiftRightLogical(dest uint, reg uint, shiftreg uint)
	shiftRightArithmetic(dest uint, reg uint, shiftreg uint)
	branchEqual()
	branchNotEqual()
	branchLessThan()
	branchLessThanUnsigned()
	branchGreaterThanOrEqual()
	branchGreaterThanOrEqualUnsigned()
	jumpAndLink(dest uint, pcOffset uint32)
	jumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32)
	loadWord(dest uint, reg uint, offset uint32)
	loadHalfWord(dest uint, reg uint, offset uint32)
	loadHalfWordUnsigned(dest uint, reg uint, offset uint32)
	loadByte(dest uint, reg uint, offset uint32)
	loadByteUnsigned(dest uint, reg uint, offset uint32)
	storeWord(src uint, reg uint, offset uint32)
	storeHalfWord(src uint, reg uint, offset uint32)
	storeByte(src uint, reg uint, offset uint32)
	csrReadAndWrite()
	csrReadAndSet()
	csrReadAndClear()
	csrReadAndWriteImmediate()
	csrReadAndSetImmediate()
	csrReadAndClearImmediate()
	envCall()
	envBreak()
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
		executor = factory.produceI(result, factory.executor)
	case U:
		executor = factory.produceU(result, factory.executor)
	case R:
		executor = factory.produceR(result, factory.executor)
	case J:
		executor = factory.produceJ(result, factory.executor)
	case B:
		executor = factory.produceB(result, factory.executor)
	case S:
		executor = factory.produceS(result, factory.executor)
	default:
		panic(fmt.Sprintf("unrecognized instruction type: %d", result.instructionType))
	}

	return executor
}

type executorI struct {
	executor riscVExecutor
	result   riscVBinaryParseResult
}

type validOperationI uint

/*These constants represent the possible immediate arithmetic operations
that are available.
*/
const (
	AddI validOperationI = iota + 1
	SLTI
	SLTIU
	AndI
	OrI
	XorI
	ShiftLeftLI
	ShiftRightLI
	ShiftRightAI
)

type executionFunctionI func(ex riscVExecutor, dest uint, reg uint, immediate uint32)

func (ex *executorI) execute() {
	immediate := uint32(ex.result.twelveBitImmediate)
	dest := uint(ex.result.fiveBitDestination)
	func3 := validOperationI(ex.result.funct3)
	src := uint(ex.result.fiveBitRegister1)

	decision := map[opCode](map[validOperationI](executionFunctionI)){
		ImmArith: map[validOperationI](executionFunctionI){
			AddI:         (riscVExecutor).addImmediate,
			SLTI:         (riscVExecutor).setLessThanImmediate,
			SLTIU:        (riscVExecutor).setLessThanImmediateUnsigned,
			AndI:         (riscVExecutor).andImmmediate,
			OrI:          (riscVExecutor).orImmediate,
			XorI:         (riscVExecutor).xorImmediate,
			ShiftLeftLI:  (riscVExecutor).shiftLeftLogicalImmediate,
			ShiftRightLI: (riscVExecutor).shiftRightLogicalImmediate,
			ShiftRightAI: (riscVExecutor).shiftRightArithmeticImmediate,
		},
		LUI: map[validOperationI](executionFunctionI){},
	}

	if m, ok := decision[ex.result.opCode]; ok {
		if f, ok := m[func3]; ok {
			f(ex.executor, dest, src, immediate)
		} else {
			panic(fmt.Sprintf("executionFunctionI: %d operation not found", func3))
		}
	} else {
		panic(fmt.Sprintf("executionFunctionI: %d opcode not found", ex.result.opCode))
	}
}
func (factory *riscVBinaryInstructionExecutionFactory) produceI(result riscVBinaryParseResult, ex riscVExecutor) binaryExecutor {

	var executor = executorI{
		executor: ex,
		result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorU struct {
	executor riscVExecutor
	result   riscVBinaryParseResult
}

func (ex *executorU) execute() {
	immediate := uint32(ex.result.twentyBitImmediate)
	dest := uint(ex.result.fiveBitDestination)

	switch ex.result.opCode {
	case LUI:
		ex.executor.loadUpperImmediate(dest, immediate)
	case AUIPC:
		ex.executor.loadUpperImmediate(dest, immediate)
	default:
		panic(fmt.Sprintf("executionFunctionU: %d opcode not found", ex.result.opCode)
	}
}
func (factory *riscVBinaryInstructionExecutionFactory) produceU(result riscVBinaryParseResult, ex riscVExecutor) binaryExecutor {

	var executor = executorU{
		executor: ex,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorR struct {
	executor riscVExecutor
}

func (ex *executorR) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceR(result riscVBinaryParseResult, ex riscVExecutor) binaryExecutor {

	var executor = executorR{
		executor: ex,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorJ struct {
	executor riscVExecutor
}

func (ex *executorJ) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceJ(result riscVBinaryParseResult, ex riscVExecutor) binaryExecutor {

	var executor = executorJ{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorB struct {
	executor riscVExecutor
}

func (ex *executorB) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceB(result riscVBinaryParseResult, ex riscVExecutor) binaryExecutor {

	var executor = executorB{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

type executorS struct {
	executor riscVExecutor
}

func (ex *executorS) execute() {

}
func (factory *riscVBinaryInstructionExecutionFactory) produceS(result riscVBinaryParseResult, ex riscVExecutor) binaryExecutor {

	var executor = executorS{}
	// This works because Go does Pointer Escape analysis.
	return &executor
}
