package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExecutorMemoryMock struct {
	mock.Mock
	val uint32
}

type InstructionExecutorSuite struct {
	suite.Suite
	executor RiscVInstructionExecutor
	memory   *ExecutorMemoryMock
}

func (m *ExecutorMemoryMock) Get(address uint32) uint32 {
	return m.val
}

func (m *ExecutorMemoryMock) Set(address uint32, val uint32, bits uint) {
	m.val = val
}

func TestInstructionExecutorSuite(t *testing.T) {
	suite.Run(t, new(InstructionExecutorSuite))
}

func (suite *InstructionExecutorSuite) SetupTest() {
	operator := makeAdaptedOperator([16]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	suite.executor = RiscVInstructionExecutor{
		operator: &operator,
	}

	memory := ExecutorMemoryMock{
		val: 0,
	}
	suite.memory = &memory
}

func (suite *InstructionExecutorSuite) TestLoadWord() {
	suite.memory.val = 15
	assert := assert.New(suite.T())
	suite.executor.LoadWord(0, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", 0, 13)

	assert.Equal(uint32(15), suite.executor.get(0))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWord() {

}

func (suite *InstructionExecutorSuite) TestLoadHalfWordUnsigned() {

}

func (suite *InstructionExecutorSuite) TestByteWord() {

}

func (suite *InstructionExecutorSuite) TestLoadByteUnsigned() {

}

func (suite *InstructionExecutorSuite) TestStoreWord() {

}

func (suite *InstructionExecutorSuite) TestStoreHalfWord() {

}

func (suite *InstructionExecutorSuite) TestStoreByte() {

}

func (suite *InstructionExecutorSuite) TestAddImmediate() {
	//assert := assert.New(suite.T())

	suite.memory.val = 1 << 11
	//suite.executor.LoadWord(2, )

	//suite.executor.AddImmediate()
}
