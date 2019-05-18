package binaryInstructionExecution

import (
	"math"
	"testing"

	//"github.com/stretchr/testify/assert"
	Utils "github.com/chenhowa/os/computer/binaryInstructionExecution/bitUtils"
	Producer "github.com/chenhowa/os/computer/binaryInstructionExecution/executionFactoryProducers"
	Parser "github.com/chenhowa/os/computer/binaryInstructionExecution/instructionParsing"
	"github.com/stretchr/testify/suite"
)

type ExecutionFactorySuite struct {
	suite.Suite
	executorMock *Producer.RiscVExecutorMock
	builder      *Parser.BinaryInstructionBuilder
	factory      *RiscVBinaryInstructionExecutionFactory
}

func TestExecutionFactorySuite(t *testing.T) {
	suite.Run(t, new(ExecutionFactorySuite))
}

func (suite *ExecutionFactorySuite) SetupTest() {
	suite.executorMock = new(Producer.RiscVExecutorMock)
	builder := Parser.MakeInstructionBuilder(32)
	suite.builder = &builder
	factory := MakeRiscVInstructionExecutionFactory(suite.executorMock)
	suite.factory = &factory
}

func (suite *ExecutionFactorySuite) TestInstruction_I_XorImmediate() {
	builder := suite.builder
	builder.AddNextXBits(7, uint(Parser.ImmArith))
	builder.AddNextXBits(5, 15)
	builder.AddNextXBits(3, uint(Producer.XorI))
	builder.AddNextXBits(5, 12)
	builder.AddNextXBits(12, uint(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 2, 11)))
	instruction := builder.Build()

	suite.executorMock.On("xorImmediate", uint(15), uint(12), Utils.KeepBitsInInclusiveRange(math.MaxUint32, 2, 11))
	executor := suite.factory.Produce(uint32(instruction))
	executor.Execute()
	suite.executorMock.AssertCalled(suite.T(), "xorImmediate", uint(15), uint(12), Utils.KeepBitsInInclusiveRange(math.MaxUint32, 2, 11))
}

func (suite *ExecutionFactorySuite) TestInstruction_U_AUIPC() {

}
