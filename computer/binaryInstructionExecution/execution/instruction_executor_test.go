package execution

import (
	"math"
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
	m.Called(address)

	return m.val
}

func (m *ExecutorMemoryMock) Set(address uint32, val uint32, bits uint) {
	m.Called(address, val, bits)
	m.val = val
}

func TestInstructionExecutorSuite(t *testing.T) {
	suite.Run(t, new(InstructionExecutorSuite))
}

func (suite *InstructionExecutorSuite) SetupTest() {
	operator := makeAdaptedOperator([32]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17})
	suite.executor = RiscVInstructionExecutor{
		operator: &operator,
	}

	memory := ExecutorMemoryMock{
		val: 0,
	}
	suite.memory = &memory
}

func (suite *InstructionExecutorSuite) TestLoadWord_Basic() {
	suite.memory.val = 15
	assert := assert.New(suite.T())

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadWord(0, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	assert.Equal(uint32(15), suite.executor.get(0))

}

func (suite *InstructionExecutorSuite) TestLoadWord_Advanced() {
	suite.memory.val = 15
	assert := assert.New(suite.T())

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadWord(0, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	assert.Equal(uint32(15), suite.executor.get(0))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadWord(0, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	assert.Equal(uint32(16), suite.executor.get(0))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWord() {
	suite.memory.val = 14
	assert := assert.New(suite.T())

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadHalfWord(0, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	assert.Equal(uint32(14), suite.executor.get(0))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadHalfWord(0, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	assert.Equal(uint32(15), suite.executor.get(0))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadHalfWord(0, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	assert.Equal(uint32(16), suite.executor.get(0))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWordUnsigned() {
	suite.memory.val = 14
	assert := assert.New(suite.T())

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadHalfWordUnsigned(0, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	assert.Equal(uint32(14), suite.executor.get(0))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadHalfWordUnsigned(0, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	assert.Equal(uint32(15), suite.executor.get(0))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadHalfWordUnsigned(0, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	assert.Equal(uint32(16), suite.executor.get(0))
}

func (suite *InstructionExecutorSuite) TestLoadByte() {
	suite.memory.val = 14
	assert := assert.New(suite.T())

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadByte(0, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	assert.Equal(uint32(14), suite.executor.get(0))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadByte(0, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	assert.Equal(uint32(15), suite.executor.get(0))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadByte(0, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	assert.Equal(uint32(16), suite.executor.get(0))
}

func (suite *InstructionExecutorSuite) TestLoadByteUnsigned() {
	suite.memory.val = 14
	assert := assert.New(suite.T())

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadByteUnsigned(0, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	assert.Equal(uint32(14), suite.executor.get(0))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadByteUnsigned(0, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	assert.Equal(uint32(15), suite.executor.get(0))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadByteUnsigned(0, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	assert.Equal(uint32(16), suite.executor.get(0))

}

func (suite *InstructionExecutorSuite) TestStoreWord() {
	assert := assert.New(suite.T())

	//basic test of storing a word
	suite.memory.val = 0
	suite.memory.On("Set", uint32(13), uint32(14), uint(32))
	suite.executor.StoreWord(14, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(13), uint32(14), uint(32))
	assert.Equal(uint32(14), suite.memory.val)

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 0
	suite.memory.On("Set", uint32(0), uint32(15), uint(32)) // due to overflow.
	suite.executor.StoreWord(15, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(0), uint32(15), uint(32))
	assert.Equal(uint32(15), suite.memory.val)

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 0
	suite.memory.On("Set", uint32(1<<11), uint32(16), uint(32)) // due to overflow.
	suite.executor.StoreWord(16, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(1<<11), uint32(16), uint(32))
	assert.Equal(uint32(16), suite.memory.val)
}

func (suite *InstructionExecutorSuite) TestStoreHalfWord() {
	assert := assert.New(suite.T())

	//basic test of storing a word
	suite.memory.val = 0
	suite.memory.On("Set", uint32(13), uint32(14), uint(32))
	suite.executor.StoreHalfWord(14, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(13), uint32(14), uint(32))
	assert.Equal(uint32(14), suite.memory.val)

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 0
	suite.memory.On("Set", uint32(0), uint32(15), uint(32)) // due to overflow.
	suite.executor.StoreHalfWord(15, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(0), uint32(15), uint(32))
	assert.Equal(uint32(15), suite.memory.val)

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 0
	suite.memory.On("Set", uint32(1<<11), uint32(16), uint(32)) // due to overflow.
	suite.executor.StoreHalfWord(16, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(1<<11), uint32(16), uint(32))
	assert.Equal(uint32(16), suite.memory.val)
}

func (suite *InstructionExecutorSuite) TestStoreByte() {
	assert := assert.New(suite.T())

	//basic test of storing a word
	suite.memory.val = 0
	suite.memory.On("Set", uint32(13), uint32(14), uint(32))
	suite.executor.StoreByte(14, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(13), uint32(14), uint(32))
	assert.Equal(uint32(14), suite.memory.val)

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 0
	suite.memory.On("Set", uint32(0), uint32(15), uint(32)) // due to overflow.
	suite.executor.StoreByte(15, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(0), uint32(15), uint(32))
	assert.Equal(uint32(15), suite.memory.val)

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 0
	suite.memory.On("Set", uint32(1<<11), uint32(16), uint(32)) // due to overflow.
	suite.executor.StoreByte(16, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Set", uint32(1<<11), uint32(16), uint(32))
	assert.Equal(uint32(16), suite.memory.val)
}

func (suite *InstructionExecutorSuite) TestAddImmediate() {
	//assert := assert.New(suite.T())

	suite.memory.val = 1 << 11
	//suite.executor.LoadWord(2, )

	//suite.executor.AddImmediate()
}
