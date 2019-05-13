package execution

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	Util "github.com/chenhowa/os/computer/binaryInstructionExecution/bitUtils"
)

const resultRegister = 30

type InstructionExecutorSuite struct {
	suite.Suite
	executor RiscVInstructionExecutor
	memory   *ExecutorMemoryMock
	manager  *ExecutorInstructionManagerMock
	csr      *ExecutorCsrManagerMock
}

type ExecutorCsrManagerMock struct {
	mock.Mock
	val uint32
}

func (m *ExecutorCsrManagerMock) get(reg uint) uint32 {
	m.Called(reg)
	return m.val
}

func (m *ExecutorCsrManagerMock) set(reg uint, val uint32) {
	m.Called(reg, val)
}

type ExecutorInstructionManagerMock struct {
	mock.Mock
	pcAddress uint32
}

func (im *ExecutorInstructionManagerMock) getCurrentInstructionAddress() uint32 {
	im.Called()
	return im.pcAddress
}

func (im *ExecutorInstructionManagerMock) getNextInstructionAddress() uint32 {
	im.Called()
	return im.pcAddress + 4
}

func (im *ExecutorInstructionManagerMock) addOffsetForNextInstructionAddress(offset uint32) {
	im.Called(offset)
	im.pcAddress += offset
}

func (im *ExecutorInstructionManagerMock) loadAsNextInstructionAddress(newAddress uint32) {
	im.Called(newAddress)
	im.pcAddress = newAddress
}

type ExecutorMemoryMock struct {
	mock.Mock
	val uint32
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

	manager := ExecutorInstructionManagerMock{
		pcAddress: 0,
	}
	suite.manager = &manager

	csr := ExecutorCsrManagerMock{
		val: 22,
	}
	suite.csr = &csr
}

func (suite *InstructionExecutorSuite) assertRegisterEquals(register uint, expected uint32) {
	assert.Equal(suite.T(), expected, suite.executor.Get(register))
}

func (suite *InstructionExecutorSuite) assertManagerAddressEquals(expected uint32) {
	assert.Equal(suite.T(), expected, suite.manager.pcAddress)
}

func (suite *InstructionExecutorSuite) TestLoadWord_Basic() {
	suite.memory.val = 15

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadWord(resultRegister, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(resultRegister, uint32(15))
}

func (suite *InstructionExecutorSuite) TestLoadWord_Advanced() {
	suite.memory.val = 15

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadWord(resultRegister, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(resultRegister, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadWord(resultRegister, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(resultRegister, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWord() {
	suite.memory.val = 14

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadHalfWord(resultRegister, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(resultRegister, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadHalfWord(resultRegister, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(resultRegister, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadHalfWord(resultRegister, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(resultRegister, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadHalfWordUnsigned() {
	suite.memory.val = 14

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadHalfWordUnsigned(resultRegister, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(resultRegister, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadHalfWordUnsigned(resultRegister, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(resultRegister, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadHalfWordUnsigned(resultRegister, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(resultRegister, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadByte() {

	// basic test of loading a word
	suite.memory.val = 14
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadByte(resultRegister, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(resultRegister, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadByte(resultRegister, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(resultRegister, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadByte(resultRegister, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(resultRegister, uint32(16))
}

func (suite *InstructionExecutorSuite) TestLoadByteUnsigned() {
	suite.memory.val = 14

	// basic test of loading a word
	suite.memory.On("Get", uint32(13))
	suite.executor.LoadByteUnsigned(resultRegister, 1, 12, suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(13))
	suite.assertRegisterEquals(resultRegister, uint32(14))

	//test 12-bit sign extension is occurring if 12th bit is 1
	suite.memory.val = 15
	suite.memory.On("Get", uint32(0)).Return() // duet to overflow.
	suite.executor.LoadByteUnsigned(resultRegister, 1, uint32(math.MaxUint32), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(0))
	suite.assertRegisterEquals(resultRegister, uint32(15))

	//test 12-bit sign extension is not occurring if 12th bit is 0
	suite.memory.val = 16
	suite.memory.On("Get", uint32(1<<11)).Return() // due to overflow.
	suite.executor.LoadByteUnsigned(resultRegister, 1, uint32(math.MaxUint32-(1<<11)), suite.memory)
	suite.memory.AssertCalled(suite.T(), "Get", uint32(1<<11))
	suite.assertRegisterEquals(resultRegister, uint32(16))
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
	suite.executor.AddImmediate(resultRegister, 1, 1<<11)
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 31)+1)

	// Test that sign-extension fo 12th bit works when it is 0
	suite.executor.AddImmediate(resultRegister, 1, 1<<10)
	suite.assertRegisterEquals(resultRegister, uint32((1<<10)+1))
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
	suite.executor.SetLessThanImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 11)) // compared to -1 in two's complement after sign extension
	suite.assertRegisterEquals(resultRegister, 1)

	// Test that sign extension does not occur when 12th bit is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.assertRegisterEquals(1, math.MaxUint32-1)                                                              // -1 in two's complement
	suite.executor.SetLessThanImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)) // this is some positive number, since no sign extension
	suite.assertRegisterEquals(resultRegister, 1)
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
	suite.executor.SetLessThanImmediateUnsigned(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 11))
	suite.assertRegisterEquals(resultRegister, 1)

	// Test that sign extension does not occur when 12th bit is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.assertRegisterEquals(1, math.MaxUint32-1)
	suite.executor.SetLessThanImmediateUnsigned(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10))
	suite.assertRegisterEquals(resultRegister, 0)
}

func (suite *InstructionExecutorSuite) TestAndImmediate() {
	// Basic test
	suite.executor.AndImmediate(resultRegister, 15, 30)
	suite.assertRegisterEquals(resultRegister, 14)

	// Test with 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.AndImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 11))
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31))

	// Test without 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.AndImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10))
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 10))
}

func (suite *InstructionExecutorSuite) TestOrImmediate() {
	// Basic test
	suite.executor.OrImmediate(resultRegister, 10, 5) // 0b1010 OR 0b0101 == 0b1111
	suite.assertRegisterEquals(resultRegister, 15)

	// Test with 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.OrImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11))
	suite.assertRegisterEquals(resultRegister, math.MaxUint32)

	// Test without 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 9)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.OrImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 10))
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10))
}

func (suite *InstructionExecutorSuite) TestXorImmediate() {
	// Basic test
	suite.executor.XorImmediate(resultRegister, 10, 7) // 0b1010 XOR 0b0111 == 0b1101
	suite.assertRegisterEquals(resultRegister, 13)

	// Test with 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.XorImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 11))
	suite.assertRegisterEquals(resultRegister, ^Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 10))

	// Test without 11th bit
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 10)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.XorImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 10, 10))
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 9))
}

func (suite *InstructionExecutorSuite) TestLeftShiftLogicalImmediate() {
	//Basic Test
	suite.executor.ShiftLeftLogicalImmediate(resultRegister, 1, 5)
	suite.assertRegisterEquals(resultRegister, 1<<5)

	// Test with number that uses upper 27 bits as well as lower 5 bits
	suite.executor.ShiftLeftLogicalImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(resultRegister, 1<<16)
}

func (suite *InstructionExecutorSuite) TestRightShiftLogicalImmediate() {
	//Basic Test
	suite.executor.ShiftRightLogicalImmediate(resultRegister, 15, 2)
	suite.assertRegisterEquals(resultRegister, 15>>2)

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 30)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightLogicalImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 14))

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 17, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightLogicalImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 15))
}

func (suite *InstructionExecutorSuite) TestRightShiftArithmeticImmediate() {
	//Basic Test
	suite.executor.ShiftRightArithmeticImmediate(resultRegister, 15, 2)
	suite.assertRegisterEquals(resultRegister, 15>>2)

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 30)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightArithmeticImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 14))

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 17, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.ShiftRightArithmeticImmediate(resultRegister, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5)) // 16 + 32
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31))
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
	suite.executor.Add(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, 3)

	//Advanced test that checks overflow
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(1)
	suite.executor.Add(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, 1)
}

func (suite *InstructionExecutorSuite) TestSub() {
	//Basic Test
	suite.executor.Sub(resultRegister, 3, 1)
	suite.assertRegisterEquals(resultRegister, 2)

	//Advanced test that checks overflow
	suite.executor.Sub(resultRegister, 0, 1)
	suite.assertRegisterEquals(resultRegister, math.MaxUint32)
}

func (suite *InstructionExecutorSuite) TestSetLessThan() {
	//Basic Test
	suite.executor.SetLessThan(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, 1)
	suite.executor.SetLessThan(resultRegister, 2, 1)
	suite.assertRegisterEquals(resultRegister, 0)

	//Check that signed comparison is occuring
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(5)
	suite.executor.SetLessThan(resultRegister, 5, 2)
	suite.assertRegisterEquals(resultRegister, 1) // -1 is less than 5
	suite.memory.val = math.MaxUint32 - 1
	suite.LoadMemoryIntoRegisterX(6)
	suite.executor.SetLessThan(resultRegister, 5, 6)
	suite.assertRegisterEquals(resultRegister, 0) // -1 is greater than -2
}

func (suite *InstructionExecutorSuite) TestSetLessThanUnsigned() {
	//Basic Test
	suite.executor.SetLessThanUnsigned(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, 1)
	suite.executor.SetLessThanUnsigned(resultRegister, 2, 1)
	suite.assertRegisterEquals(resultRegister, 0)

	//Check that unsigned comparison is occuring
	suite.memory.val = math.MaxUint32 - 1
	suite.LoadMemoryIntoRegisterX(6)
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(5)
	suite.executor.SetLessThanUnsigned(resultRegister, 6, 5)
	suite.assertRegisterEquals(resultRegister, 1) // (math.MaxUInt32 - 1) is less than math.MaxUint32
	suite.LoadMemoryIntoRegisterX(5)
	suite.executor.SetLessThanUnsigned(resultRegister, 5, 2)
	suite.assertRegisterEquals(resultRegister, 0) // math.MaxUint32 is greater than 5
}

func (suite *InstructionExecutorSuite) TestAnd() {
	//Basic Test
	suite.executor.And(resultRegister, 10, 7)
	suite.assertRegisterEquals(resultRegister, 2)

	//advanced test
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 14, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 20)
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.And(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 14, 20))
}

func (suite *InstructionExecutorSuite) TestOr() {
	//Basic test
	suite.executor.Or(resultRegister, 10, 5)
	suite.assertRegisterEquals(resultRegister, 15)

	//advanced test
	val := Util.KeepBitsInInclusiveRange(math.MaxUint32, 12, 16)
	suite.memory.val = val
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = ^val
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.Or(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, math.MaxUint32)
}

func (suite *InstructionExecutorSuite) TestXor() {
	//Basic test
	suite.executor.Xor(resultRegister, 7, 14) // 0b0111 XOR 0b1110 = 1001
	suite.assertRegisterEquals(resultRegister, 9)

	//Advanced test
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 29)
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 29, 30)
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.Xor(resultRegister, 1, 2)
	expected := Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 28) + Util.KeepBitsInInclusiveRange(math.MaxUint32, 30, 30)
	suite.assertRegisterEquals(resultRegister, expected)
}

func (suite *InstructionExecutorSuite) TestLeftShiftLogical() {
	//Basic Test
	suite.executor.ShiftLeftLogical(resultRegister, 1, 5)
	suite.assertRegisterEquals(resultRegister, 1<<5)

	// Test with number that uses upper 27 bits as well as lower 5 bits
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5) // 16 + 32
	suite.LoadMemoryIntoRegisterX(5)
	suite.executor.ShiftLeftLogical(resultRegister, 1, 5)
	suite.assertRegisterEquals(resultRegister, 1<<16)
}

func (suite *InstructionExecutorSuite) TestRightShiftLogical() {
	//Basic Test
	suite.executor.ShiftRightLogical(resultRegister, 15, 2)
	suite.assertRegisterEquals(resultRegister, 15>>2)

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 30)
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5) // 16 + 32
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.ShiftRightLogical(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 14))

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 17, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5) // 16 + 32
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.ShiftRightLogical(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 15))
}

func (suite *InstructionExecutorSuite) TestRightShiftArithmetic() {
	//Basic Test
	suite.executor.ShiftRightArithmetic(resultRegister, 15, 2)
	suite.assertRegisterEquals(resultRegister, 15>>2)

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 0
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 30)
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5) // 16 + 32
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.ShiftRightArithmetic(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 14))

	// Test with number that uses upper 27 bits as well as lower 5 bits, where MSB is 1
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 17, 31)
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = Util.KeepBitsInInclusiveRange(math.MaxUint32, 4, 5) // 16 + 32
	suite.LoadMemoryIntoRegisterX(2)
	suite.executor.ShiftRightArithmetic(resultRegister, 1, 2)
	suite.assertRegisterEquals(resultRegister, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 31))
}

func (suite *InstructionExecutorSuite) TestBranchEqual() {
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)

	// Test if not equal
	suite.executor.BranchEqual(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	//Test if equal
	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchEqual(2, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))
}

func (suite *InstructionExecutorSuite) TestBranchNotEqual() {
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)

	//Test if equal
	suite.executor.BranchNotEqual(2, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	// Test if not equal
	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchNotEqual(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))

}

func (suite *InstructionExecutorSuite) TestBranchLessThan_Basic() {
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)
	//Test basic
	suite.executor.BranchLessThan(2, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.executor.BranchLessThan(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchLessThan(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))

}

func (suite *InstructionExecutorSuite) TestBranchLessThan_Advanced() {

	//Test signed comparison is occurring
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = math.MaxUint32 - 1
	suite.LoadMemoryIntoRegisterX(2)

	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)

	suite.executor.BranchLessThan(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.executor.BranchLessThan(3, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchLessThan(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))
}

func (suite *InstructionExecutorSuite) TestBranchLessThanUnsigned_Basic() {
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)
	//Test basic
	suite.executor.BranchLessThanUnsigned(2, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.executor.BranchLessThanUnsigned(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchLessThanUnsigned(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))

}

func (suite *InstructionExecutorSuite) TestBranchLessThanUnsigned_Advanced() {

	//Test unsigned comparison is occurring
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(1)
	suite.memory.val = math.MaxUint32 - 1
	suite.LoadMemoryIntoRegisterX(2)

	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)

	suite.executor.BranchLessThanUnsigned(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.executor.BranchLessThanUnsigned(1, 3, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchLessThanUnsigned(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))
}

func (suite *InstructionExecutorSuite) TestBranchGreaterThanOrEqual_Basic() {
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)
	//Test basic
	suite.executor.BranchGreaterThanOrEqual(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchGreaterThanOrEqual(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual+1))
	suite.executor.BranchGreaterThanOrEqual(2, 2, offset+1, suite.manager)
	suite.assertManagerAddressEquals(actual + actual + 1)

}

func (suite *InstructionExecutorSuite) TestBranchGreaterThanOrEqual_Advanced() {

	//Test signed comparison is occurring
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(1)      //register 1 contains -1
	suite.memory.val = math.MaxUint32 - 1 // register 2 contains -2
	suite.LoadMemoryIntoRegisterX(2)

	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)

	suite.executor.BranchGreaterThanOrEqual(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.executor.BranchGreaterThanOrEqual(1, 3, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchGreaterThanOrEqual(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))
}

func (suite *InstructionExecutorSuite) TestBranchGreaterThanOrEqualUnsigned_Basic() {
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)
	//Test basic
	suite.executor.BranchGreaterThanOrEqualUnsigned(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchGreaterThanOrEqualUnsigned(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual+1))
	suite.executor.BranchGreaterThanOrEqualUnsigned(2, 2, offset+1, suite.manager)
	suite.assertManagerAddressEquals(actual + actual + 1)

}

func (suite *InstructionExecutorSuite) TestBranchGreaterThanOrEqualUnsigned_Advanced() {

	//Test signed comparison is occurring
	suite.memory.val = math.MaxUint32
	suite.LoadMemoryIntoRegisterX(1)      //register 1 contains -1
	suite.memory.val = math.MaxUint32 - 1 // register 2 contains -2
	suite.LoadMemoryIntoRegisterX(2)

	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 12)
	actual := Util.KeepBitsInInclusiveRange(math.MaxUint32, 11, 11)

	suite.executor.BranchGreaterThanOrEqualUnsigned(2, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.executor.BranchGreaterThanOrEqualUnsigned(3, 1, offset, suite.manager)
	suite.assertManagerAddressEquals(0)
	suite.manager.AssertNotCalled(suite.T(), "addOffsetForNextInstructionAddress", mock.Anything)

	suite.manager.On("addOffsetForNextInstructionAddress", uint32(actual))
	suite.executor.BranchGreaterThanOrEqualUnsigned(1, 2, offset, suite.manager)
	suite.assertManagerAddressEquals(actual)
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(actual))
}

func (suite *InstructionExecutorSuite) TestJumpAndLink() {
	suite.manager.pcAddress = 1
	suite.manager.On("addOffsetForNextInstructionAddress", uint32(math.MaxUint32))
	suite.manager.On("getNextInstructionAddress")
	suite.executor.JumpAndLink(1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 0, 19), suite.manager)
	suite.assertRegisterEquals(1, 1+4)
	suite.assertManagerAddressEquals(0) // 1 + math.MaxUint32
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", uint32(math.MaxUint32))
	suite.manager.AssertCalled(suite.T(), "getNextInstructionAddress")

	suite.manager.pcAddress = 1
	offset := Util.KeepBitsInInclusiveRange(math.MaxUint32, 18, 18)
	suite.manager.On("addOffsetForNextInstructionAddress", offset)
	suite.executor.JumpAndLink(1, offset, suite.manager)
	suite.assertRegisterEquals(1, 1+4)
	suite.assertManagerAddressEquals(offset + 1) // 1 + math.MaxUint32
	suite.manager.AssertCalled(suite.T(), "addOffsetForNextInstructionAddress", offset)
}

func (suite *InstructionExecutorSuite) TestJumpAndLinkRegister() {
	suite.manager.pcAddress = 1
	suite.manager.On("getNextInstructionAddress")
	suite.manager.On("loadAsNextInstructionAddress", uint32(math.MaxUint32-1))
	suite.executor.JumpAndLinkRegister(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 11), suite.manager)
	suite.assertRegisterEquals(30, 1+4)
	suite.assertManagerAddressEquals(math.MaxUint32 - 1)

	suite.manager.pcAddress = 2
	suite.manager.On("getNextInstructionAddress")
	suite.manager.On("loadAsNextInstructionAddress", Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 10))
	suite.executor.JumpAndLinkRegister(30, 1, Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 10), suite.manager)
	suite.assertRegisterEquals(30, 2+4)
	suite.assertManagerAddressEquals(Util.KeepBitsInInclusiveRange(math.MaxUint32, 1, 10))
}

func (suite *InstructionExecutorSuite) TestCsrReadAndWrite() {
	suite.csr.val = 15
	suite.csr.On("get", mock.Anything)
	suite.csr.On("set", uint(3), uint32(2))
	suite.executor.CsrReadAndWrite(0, 2, 3, suite.csr)
	suite.csr.AssertNotCalled(suite.T(), "get", mock.Anything)
	suite.csr.AssertCalled(suite.T(), "set", uint(3), uint32(2))
	suite.assertRegisterEquals(2, 2)
	suite.assertRegisterEquals(0, 0)

	suite.csr.val = 16
	suite.csr.On("get", uint(4))
	suite.csr.On("set", uint(4), uint32(2))
	suite.executor.CsrReadAndWrite(1, 2, 4, suite.csr)
	suite.csr.AssertCalled(suite.T(), "get", uint(4))
	suite.csr.AssertCalled(suite.T(), "set", uint(4), uint32(2))
	suite.assertRegisterEquals(1, 16)
	suite.assertRegisterEquals(2, 2)
}

func (suite *InstructionExecutorSuite) TestCsrReadAndSet() {

}

func (suite *InstructionExecutorSuite) TestCsrReadAndClear() {

}

func (suite *InstructionExecutorSuite) TestCsrReadAndWriteImmediate() {

}

func (suite *InstructionExecutorSuite) TestCsrReadAndSetImmediate() {

}

func (suite *InstructionExecutorSuite) TestCsrReadAndClearImmediate() {

}
