package executionFactoryProducers

import (
	"fmt"

	Parser "github.com/chenhowa/computer/lib/binaryInstructionExecution/instructionParsing"
)

/*ExecutorB contains instructions for executing
a B-type instruction ParseResult
*/
type ExecutorB struct {
	Executor RiscVExecutor
	Result   Parser.RiscVBinaryParseResult
}

type executionFunctionB func(ex RiscVExecutor, src1 uint, src2 uint, offset uint32)

type validOperationB uint

/*These constants specify valid operations for
B-type instructions
*/
const (
	Beq validOperationB = iota + 1
	Bneq
	Blt
	Bltu
	Bge
	Bgeu
)

/*Execute will execute the B-type instruction
 */
func (ex *ExecutorB) Execute() {
	immediate := uint32(ex.Result.TwelveBitImmediate)
	src1 := uint(ex.Result.FiveBitRegister1)
	src2 := uint(ex.Result.FiveBitRegister2)
	func3 := validOperationB(ex.Result.Funct3)

	decision := map[Parser.OpCode](map[validOperationB](executionFunctionB)){
		Parser.Branch: map[validOperationB](executionFunctionB){
			Beq:  (RiscVExecutor).branchEqual,
			Bneq: (RiscVExecutor).branchNotEqual,
			Blt:  (RiscVExecutor).branchLessThan,
			Bltu: (RiscVExecutor).branchLessThanUnsigned,
			Bge:  (RiscVExecutor).branchGreaterThanOrEqual,
			Bgeu: (RiscVExecutor).branchGreaterThanOrEqualUnsigned,
		},
	}

	if m, ok := decision[ex.Result.OpCode]; ok {
		if f, ok := m[func3]; ok {
			f(ex.Executor, src1, src2, immediate)
		} else {
			panic(fmt.Sprintf("executionFunctionB: %d operation not found", func3))
		}
	} else {
		panic(fmt.Sprintf("executionFunctionB: %d opcode not found", ex.Result.OpCode))
	}
}
