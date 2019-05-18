package binaryInstructionExecution

import (
	"math"
	"testing"

	//"github.com/stretchr/testify/assert"
	Utils "github.com/chenhowa/computer/lib/binaryInstructionExecution/bitUtils"
	Producer "github.com/chenhowa/computer/lib/binaryInstructionExecution/executionFactoryProducers"
	Parser "github.com/chenhowa/computer/lib/binaryInstructionExecution/instructionParsing"
	"github.com/stretchr/testify/suite"
)

type ExecutionFactorySuite struct {
	suite.Suite
	executorMock *Producer.RiscVExecutorMock
	factory      *RiscVBinaryInstructionExecutionFactory
}

func TestExecutionFactorySuite(t *testing.T) {
	suite.Run(t, new(ExecutionFactorySuite))
}

func (suite *ExecutionFactorySuite) SetupTest() {
	suite.executorMock = new(Producer.RiscVExecutorMock)
	factory := MakeRiscVInstructionExecutionFactory(suite.executorMock)
	suite.factory = &factory
}

func (suite *ExecutionFactorySuite) TestInstruction_I_XorImmediate() {
	instruction := BuildInstructionI(uint(Parser.ImmArith), 15, uint(Producer.XorI), 12, uint(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 2, 11)))

	suite.executorMock.On("xorImmediate", uint(15), uint(12), Utils.KeepBitsInInclusiveRange(math.MaxUint32, 2, 11))
	executor := suite.factory.Produce(uint32(instruction))
	executor.Execute()
	suite.executorMock.AssertCalled(suite.T(), "xorImmediate", uint(15), uint(12), Utils.KeepBitsInInclusiveRange(math.MaxUint32, 2, 11))
}

func (suite *ExecutionFactorySuite) TestInstruction_U_AUIPC() {
	instruction := BuildInstructionU(uint(Parser.AUIPC), 20, uint(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 1, 19)))

	suite.executorMock.On("addUpperImmediateToPC", uint(20), Utils.KeepBitsInInclusiveRange(math.MaxUint32, 1, 19))
	suite.factory.Produce(uint32(instruction)).Execute()
	suite.executorMock.AssertCalled(suite.T(), "addUpperImmediateToPC", uint(20), Utils.KeepBitsInInclusiveRange(math.MaxUint32, 1, 19))

}
