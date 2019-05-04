package executionFactoryProducers

import (
	"fmt"

	Parser "github.com/chenhowa/os/computer/binaryInstructionExecution/instructionParsing"
)

/*ExecutorS contains instructions for executing
a S-type instruction ParseResult
*/
type ExecutorS struct {
	Executor RiscVExecutor
	Result   Parser.RiscVBinaryParseResult
}

type executionFunctionS func(ex RiscVExecutor, reg1 uint, reg2 uint, offset uint32)

type validOperationS uint

/* These constants represent valid operations when
instruction is S-type and OpCode is Store
*/
const (
	StoreWord validOperationS = iota
	StoreHalfWord
	StoreByte
)

/*Execute will execute the S-type instruction
 */
func (ex *ExecutorS) Execute() {
	immediate := uint32(ex.Result.TwelveBitImmediate)
	src := uint(ex.Result.FiveBitRegister2)
	base := uint(ex.Result.FiveBitRegister1)
	func3 := validOperationS(ex.Result.Funct3)

	decision := map[Parser.OpCode](map[validOperationS](executionFunctionS)){
		Parser.Store: map[validOperationS](executionFunctionS){
			StoreWord:     (RiscVExecutor).storeWord,
			StoreHalfWord: (RiscVExecutor).storeHalfWord,
			StoreByte:     (RiscVExecutor).storeByte,
		},
	}

	if m, ok := decision[ex.Result.OpCode]; ok {
		if f, ok := m[func3]; ok {
			f(ex.Executor, base, src, immediate)
		} else {
			panic(fmt.Sprintf("executionFunctionS: %d operation not found", func3))
		}
	} else {
		panic(fmt.Sprintf("executionFunctionS: %d opcode not found", ex.Result.OpCode))
	}
}
