package execution

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	Util "github.com/chenhowa/os/computer/binaryInstructionExecution/bitUtils"
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
	operator := makeAdaptedOperator([32]uint32{24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17})
	suite.executor = RiscVInstructionExecutor{
		operator: &operator,
	}

	memory := ExecutorMemoryMock{
		val: 0,
	}
	suite.memory = &memory
}

func (suite *InstructionExecutorSuite) assertRegisterEquals(register uint, expected uint32) {
	assert.Equal(suite.T(), expected, suite.executor.Get(register))
}

func (suite *InstructionExecutorSuite) TestLoadWord_Basic() {
	suite.memory.val = 15

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadWord(30, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(30, uint32(15))
}

func (suite *InstructionExecutorSuite) TestLoadWord_Advanced() {
	suite.memory.val = 15

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadWord(30, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(30, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadWord(30, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(30, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWord() {
	suite.memory.val = 14

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadHalfWord(30, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(30, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadHalfWord(30, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(30, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadHalfWord(30, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(30, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWordUnsigned() {
	suite.memory.val = 14

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadHalfWordUnsigned(30, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(30, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadHalfWordUnsigned(30, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(30, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadHalfWordUnsigned(30, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(30, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadByte() {

	// basic test of loading a word
	suite.memory.val = 14
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadByte(30, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(30, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadByte(30, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(30, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadByte(30, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(30, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadByteUnsigned() {
	suite.memory.val = 14

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadByteUnsigned(30, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(30, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadByteUnsigned(30, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(30, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadByteUnsigned(30, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(30, uint32(16))
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

func (suite *InstructionExecutorSuite) LoadMemoryIntoRegisterX(x uint) {
	suite.memory.On("Get", mock.Anything)
	suite.executor.LoadWord(x, 5, 5, suite.memory)
}

func (suite *InstructionExecutorSuite) TestAddImmediate() {

	// Test that sign-extension of 12th bit works when it is 1
	suite.executor.AddImmediate(30, 1, 1<<11)
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 31)+1)

	// Test that sign-extension fo 12th bit works when it is 0
	suite.executor.AddImmediate(30, 1, 1<<10)
	suite.assertRegisterEquals(30, uint32((1<<10)+1))
}

func (suite *InstructionExecutorSuite) TestSLTI() {
	//Basic test of setting functionality
	suite.executor.SetLessThanImmediate(3, 2, 1)
	suite.assertRegisterEquals(3, 0)
	suite.executor.SetLessThanImmediate(3, 1, 2)
	suite.assertRegisterEquals(3, 1)
	suite.executor.SetLessThanImmediate(3, 1, 1)
	suite.assertRegisterEquals(3, 0)

	// Test that the sign-extension of 12th bit works when it is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.assertRegisterEquals(1, math.MaxUint32-1)
	// this is -2 in two's complement
	// compared to -1 in two's complement after sign extension
	suite.executor.SetLessThanImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 11)) // compared to -1 in two's complement after sign extension
	suite.assertRegisterEquals(30, 1)

	// Test that sign extension does not occur when 12th bit is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.assertRegisterEquals(1, math.MaxUint32-1)                                                  // -1 in two's complement
	suite.executor.SetLessThanImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)) // this is some positive number, since no sign extension
	suite.assertRegisterEquals(30, 1)
}

func (suite *InstructionExecutorSuite) TestSLTIUnsigned() {

	//Basic test of setting functionality
	suite.executor.SetLessThanImmediateUnsigned(3, 2, 1)
	suite.assertRegisterEquals(3, 0)
	suite.executor.SetLessThanImmediateUnsigned(3, 1, 2)
	suite.assertRegisterEquals(3, 1)
	suite.executor.SetLessThanImmediateUnsigned(3, 1, 1)
	suite.assertRegisterEquals(3, 0)

	// Test that the sign-extension of 12th bit works when it is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.assertRegisterEquals(1, math.MaxUint32-1)
	suite.executor.SetLessThanImmediateUnsigned(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 11))
	suite.assertRegisterEquals(30, 1)

	// Test that sign extension does not occur when 12th bit is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.assertRegisterEquals(1, math.MaxUint32-1)
	suite.executor.SetLessThanImmediateUnsigned(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10))
	suite.assertRegisterEquals(30, 0)
}

func (suite *InstructionExecutorSuite) TestAndImmediate() {
	// Basic test
	suite.executor.AndImmediate(30, 15, 30)
	suite.assertRegisterEquals(30, 14)

	// Test with 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.AndImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 11))
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31))

	// Test without 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.AndImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10))
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 10))
}

func (suite *InstructionExecutorSuite) TestOrImmediate() {
	// Basic test
	suite.executor.OrImmediate(30, 10, 5) // 0b1010 OR 0b0101 == 0b1111
	suite.assertRegisterEquals(30, 15)

	// Test with 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.OrImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11))
	suite.assertRegisterEquals(30, math.MaxUint32)

	// Test without 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 9)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.OrImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 10))
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10))
}

func (suite *InstructionExecutorSuite) TestXorImmediate() {
	// Basic test
	suite.executor.XorImmediate(30, 10, 7) // 0b1010 XOR 0b0111 == 0b1101
	suite.assertRegisterEquals(30, 13)

	// Test with 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.XorImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 11))
	suite.assertRegisterEquals(30, ^Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 10))

	// Test without 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.XorImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 10))
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 9))
}

func (suite *InstructionExecutorSuite) TestLeftShiftLogicalImmediate() {
	//Basic Test
	suite.executor.ShiftLeftLogicalImmediate(30, 1, 5)
	suite.assertRegisterEquals(30, 1<<5)

	// Test with number that uses upper 27 bits as well as lower 5 bits
	suite.executor.ShiftLeftLogicalImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(30, 1<<16)
}

func (suite *InstructionExecutorSuite) TestRightShiftLogicalImmediate() {
	//Basic Test
	suite.executor.ShiftRightLogicalImmediate(30, 15, 2)
	suite.assertRegisterEquals(30, 15>>2)

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 30)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightLogicalImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 14))

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 17, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightLogicalImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 15))
}

func (suite *InstructionExecutorSuite) TestRightShiftArithmeticImmediate() {
	//Basic Test
	suite.executor.ShiftRightArithmeticImmediate(30, 15, 2)
	suite.assertRegisterEquals(30, 15>>2)

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 30)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightArithmeticImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 14))

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 17, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightArithmeticImmediate(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(30, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31))
}

func (suite *InstructionExecutorSuite) TestLUI() {
	suite.executor.LoadUpperImmediate(1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 25))
	suite.assertRegisterEquals(1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 13, 31))
}

func (suite *InstructionExecutorSuite) TestAUIPC() {
	suite.assertRegisterEquals(0, 5)
}

func (suite *InstructionExecutorSuite) TestAdd() {
	//Basic Test
	suite.executor.Add(30, 1, 2)
	suite.assertRegisterEquals(30, 3)

	//Advanced test that checks overflow
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.Add(30, 1, 2)
	suite.assertRegisterEquals(30, 1)
}

func (suite *InstructionExecutorSuite) TestSub() {
	//Basic Test
	suite.executor.Sub(30, 3, 1)
	suite.assertRegisterEquals(30, 2)

	//Advanced test that checks overflow
	suite.executor.Sub(30, 0, 1)
	suite.assertRegisterEquals(30, math.MaxUint32)
}

func (suite *InstructionExecutorSuite) TestSetLessThan() {

}

func (suite *InstructionExecutorSuite) TestSetLessThanUnsigned() {

}
