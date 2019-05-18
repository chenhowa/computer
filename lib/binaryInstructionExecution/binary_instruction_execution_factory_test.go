package binaryInstructionExecution

import (
	"math"
	"testing"

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

func (suite *ExecutionFactorySuite) TestInstruction_R_ShiftRightArithmetic() {
	instruction := BuildInstructionR(uint(Parser.RegArith), 11, uint(Producer.SRA), 4, 5, uint(Producer.F1))
	suite.executorMock.On("shiftRightArithmetic", uint(11), uint(4), uint(5))
	suite.factory.Produce(uint32(instruction)).Execute()
	suite.executorMock.AssertCalled(suite.T(), "shiftRightArithmetic", uint(11), uint(4), uint(5))
}

func (suite *ExecutionFactorySuite) TestInstruction_J_JAL() {
	instruction := uint32(BuildInstructionJ(uint(Parser.JAL), 15, 46))
	suite.executorMock.On("jumpAndLink", uint(15), uint32(46))
	suite.factory.Produce(instruction).Execute()
	suite.executorMock.AssertCalled(suite.T(), "jumpAndLink", uint(15), uint32(46))
}

func (suite *ExecutionFactorySuite) TestInstruction_B_BGE() {
	instruction := uint32(BuildInstructionB(uint(Parser.Branch), uint(Producer.Bge), 2, 3, 103))
	suite.executorMock.On("branchGreaterThanOrEqual", uint(2), uint(3), uint32(103))
	suite.factory.Produce(instruction).Execute()
	suite.executorMock.AssertCalled(suite.T(), "branchGreaterThanOrEqual", uint(2), uint(3), uint32(103))
}

func (suite *ExecutionFactorySuite) TestInstruction_S_StoreByte() {
	instruction := uint32(BuildInstructionS(uint(Parser.Store), uint(Producer.StoreByte), 12, 30, 1000))
	suite.executorMock.On("storeByte", uint(12), uint(30), uint32(1000))
	suite.factory.Produce(instruction).Execute()
	suite.executorMock.AssertCalled(suite.T(), "storeByte", uint(12), uint(30), uint32(1000))
}
