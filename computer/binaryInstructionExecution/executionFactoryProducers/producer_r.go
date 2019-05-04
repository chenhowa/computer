package executionFactoryProducers

import "fmt"
import Parser "github.com/chenhowa/os/computer/binaryInstructionExecution/instructionParsing"

type validOperationR uint

/*These constants define the valid operation codes for
R-type instructions
*/
const (
	Add validOperationR = iota + 1
	SLT
	SLTU
	And
	Or
	Xor
	SLL
	SRL
	Sub
	SRA
)

/*These constants are valid Funct7 constants
for an R-type instruction.
*/
const (
	F0 funct7 = iota + 1
	F1
)

type executionFunctionR func(ex RiscVExecutor, dest uint, reg1 uint, reg2 uint)

/*ExecutorR stores and canc execute instructions for
a R-type instruction ParseResult
*/
type ExecutorR struct {
	Executor RiscVExecutor
	Result   Parser.RiscVBinaryParseResult
}

type funct7 uint

func (ex *ExecutorR) Execute() {
	src1 := uint(ex.Result.FiveBitRegister1)
	src2 := uint(ex.Result.FiveBitRegister2)
	dest := uint(ex.Result.FiveBitDestination)
	func7 := funct7(ex.Result.Funct7)
	func3 := validOperationR(ex.Result.Funct3)

	decision := map[Parser.OpCode](map[funct7](map[validOperationR](executionFunctionR))){
		Parser.RegArith: map[funct7](map[validOperationR](executionFunctionR)){
			F0: map[validOperationR](executionFunctionR){
				Add:  (RiscVExecutor).add,
				SLT:  (RiscVExecutor).setLessThan,
				SLTU: (RiscVExecutor).setLessThanUnsigned,
				And:  (RiscVExecutor).and,
				Or:   (RiscVExecutor).or,
				Xor:  (RiscVExecutor).xor,
				SLL:  (RiscVExecutor).shiftLeftLogical,
				SRL:  (RiscVExecutor).shiftRightLogical,
			},
			F1: map[validOperationR](executionFunctionR){
				Sub: (RiscVExecutor).sub,
				SRA: (RiscVExecutor).shiftRightArithmetic,
			},
		},
	}

	if m1, ok := decision[ex.Result.OpCode]; ok {
		if m2, ok := m1[func7]; ok {
			if f, ok := m2[func3]; ok {
				f(ex.Executor, dest, src1, src2)
			} else {
				panic(fmt.Sprintf("executionFunctionR: %d operation not found", func3))
			}
		} else {
			panic(fmt.Sprintf("executionFunctionR: %d func7 not found", func7))
		}
	} else {
		panic(fmt.Sprintf("executionFunctionR: %d opcode not found", ex.Result.OpCode))

	}
}
