package executionFactoryProducers

import Parser "github.com/chenhowa/computer/lib/binaryInstructionExecution/instructionParsing"

import "fmt"

/*ExecutorJ contains instructions for executing
a J-type instruction ParseResult
*/
type ExecutorJ struct {
	Executor RiscVExecutor
	Result   Parser.RiscVBinaryParseResult
}

/*Execute will execute the J-type instruction
 */
func (ex *ExecutorJ) Execute() {
	immediate := uint32(ex.Result.TwentyBitImmediate)
	dest := uint(ex.Result.FiveBitDestination)

	switch ex.Result.OpCode {
	case Parser.JAL:
		ex.Executor.jumpAndLink(dest, immediate)
	default:
		panic(fmt.Sprintf("executionFunctionJ: %d opcode not found", ex.Result.OpCode))
	}
}
