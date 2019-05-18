package executionFactoryProducers

import "fmt"
import Parser "github.com/chenhowa/computer/lib/binaryInstructionExecution/instructionParsing"

/*ExecutorU stores execution of a
U-type instruction ParseResult
*/
type ExecutorU struct {
	Executor RiscVExecutor
	Result   Parser.RiscVBinaryParseResult
}

/*Execute will execute the U-type instruction
 */
func (ex *ExecutorU) Execute() {
	immediate := uint32(ex.Result.TwentyBitImmediate)
	dest := uint(ex.Result.FiveBitDestination)

	switch ex.Result.OpCode {
	case Parser.LUI:
		ex.Executor.loadUpperImmediate(dest, immediate)
	case Parser.AUIPC:
		ex.Executor.loadUpperImmediate(dest, immediate)
	default:
		panic(fmt.Sprintf("executionFunctionU: %d opcode not found", ex.Result.OpCode))
	}
}
