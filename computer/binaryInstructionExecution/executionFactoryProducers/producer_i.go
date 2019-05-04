package executionFactoryProducers

import "fmt"
import Parser "github.com/chenhowa/os/computer/binaryInstructionExecution/instructionParsing"

/*ExecutorI stores instructions for executing a given set
of commands for an I-type instruction ParseResult
*/
type ExecutorI struct {
	Executor RiscVExecutor
	Result   Parser.RiscVBinaryParseResult
}

type validOperationI uint

/*These constants represent the possible operations
that are available for I-type instructions
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
	JALR
)

type executionFunctionI func(ex RiscVExecutor, dest uint, reg uint, immediate uint32)

/*Execute will execute the I-type instruction
 */
func (ex *ExecutorI) Execute() {
	immediate := uint32(ex.Result.TwelveBitImmediate)
	dest := uint(ex.Result.FiveBitDestination)
	func3 := validOperationI(ex.Result.Funct3)
	src := uint(ex.Result.FiveBitRegister1)

	decision := map[Parser.OpCode](map[validOperationI](executionFunctionI)){
		Parser.ImmArith: map[validOperationI](executionFunctionI){
			AddI:         (RiscVExecutor).addImmediate,
			SLTI:         (RiscVExecutor).setLessThanImmediate,
			SLTIU:        (RiscVExecutor).setLessThanImmediateUnsigned,
			AndI:         (RiscVExecutor).andImmmediate,
			OrI:          (RiscVExecutor).orImmediate,
			XorI:         (RiscVExecutor).xorImmediate,
			ShiftLeftLI:  (RiscVExecutor).shiftLeftLogicalImmediate,
			ShiftRightLI: (RiscVExecutor).shiftRightLogicalImmediate,
			ShiftRightAI: (RiscVExecutor).shiftRightArithmeticImmediate,
		},
		Parser.JALR: map[validOperationI](executionFunctionI){
			JALR: (RiscVExecutor).jumpAndLinkRegister,
		},
	}

	if m, ok := decision[ex.Result.OpCode]; ok {
		if f, ok := m[func3]; ok {
			f(ex.Executor, dest, src, immediate)
		} else {
			panic(fmt.Sprintf("executionFunctionI: %d operation not found", func3))
		}
	} else {
		panic(fmt.Sprintf("executionFunctionI: %d opcode not found", ex.Result.OpCode))
	}
}
