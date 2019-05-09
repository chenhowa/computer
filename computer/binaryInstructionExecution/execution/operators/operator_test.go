package operators

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OperatorSuite struct {
	suite.Suite
	operator Operator
	memory   MemoryMock
}

type MemoryMock struct {
	val uint32
}

func (m *MemoryMock) Get(address uint16) uint32 {
	return m.val
}

func (m *MemoryMock) Set(address uint16, val uint32, bits uint) {
	m.val = val
}

func (suite *OperatorSuite) SetupTest() {
	suite.operator = Operator{
		registers: [16]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	}
	suite.memory = MemoryMock{}
}

func (suite *OperatorSuite) AssertRegisterEquals(reg uint16, val uint32) {
	assert := assert.New(suite.T())
	actual := suite.operator.registers[reg]
	assert.Equal(val, actual)
}

func TestRun(t *testing.T) {
	suite.Run(t, new(OperatorSuite))
}

func (suite *OperatorSuite) TestLoadStore() {
	var assert = assert.New(suite.T())

	var start_value uint32 = 10
	suite.memory.val = start_value
	assert.Equal(suite.memory.val, start_value)

	suite.operator.Load(0, 5, &suite.memory)
	assert.Equal(suite.operator.registers[0], start_value)

	suite.operator.Load(15, 5, &suite.memory)
	assert.Equal(suite.operator.registers[15], start_value)

	suite.operator.Store(1, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(1))
}

func (suite *OperatorSuite) TestAdd() {
	var assert = assert.New(suite.T())

	var start_value uint32 = 10
	suite.memory.val = start_value
	suite.operator.Load(0, 5, &suite.memory)
	suite.operator.Load(15, 5, &suite.memory)

	suite.operator.Add(5, 0, 15)
	assert.Equal(suite.operator.registers[5], start_value*2)
}

func (suite *OperatorSuite) TestSub() {
	var assert = assert.New(suite.T())

	suite.operator.Sub(0, 3, 1)
	suite.operator.Store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(2))
}

func (suite *OperatorSuite) TestBitAnd() {
	var assert = assert.New(suite.T())

	suite.operator.Bit_and(0, 1, 3)
	suite.operator.Store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(1))
}

func (suite *OperatorSuite) TestBitOr() {
	var assert = assert.New(suite.T())

	suite.operator.Bit_or(0, 1, 2)
	suite.operator.Store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(3))
}

func (suite *OperatorSuite) TestBitXor() {
	suite.operator.Bit_xor(0, 5, 3) // 0b101 ^ 0b011 = 6
	suite.AssertRegisterEquals(0, 6)
	suite.AssertRegisterEquals(3, 3)
	suite.AssertRegisterEquals(5, 5)
}

func (suite *OperatorSuite) TestBitNot() {
	suite.memory.val = math.MaxUint32 - 2
	suite.operator.Load(5, 3, &suite.memory)

	suite.operator.Bit_not(0, 5)
	suite.AssertRegisterEquals(0, 2)
}

func (suite *OperatorSuite) TestMultiply() {
	var assert = assert.New(suite.T())
	suite.operator.Multiply(0, 2, 3)
	suite.operator.Store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(6))
}

func (suite *OperatorSuite) TestDivide() {
	suite.operator.Divide(0, 1, 5, 3)
	suite.AssertRegisterEquals(0, uint32(1))
	suite.AssertRegisterEquals(1, uint32(2))
}

func (suite *OperatorSuite) TestAddImmediate() {
	suite.operator.Add_immediate(0, 1, 10)
	suite.AssertRegisterEquals(0, 11)
	suite.AssertRegisterEquals(1, 1)
}

func (suite *OperatorSuite) TestAndImmediate() {
	suite.operator.Bit_and_immediate(0, 2, 7)
	suite.AssertRegisterEquals(0, 2)
	suite.AssertRegisterEquals(2, 2)
}

func (suite *OperatorSuite) TestOrImmediate() {
	suite.operator.Bit_or_immediate(0, 5, 2)
	suite.AssertRegisterEquals(0, 7)
	suite.AssertRegisterEquals(5, 5)
}

func (suite *OperatorSuite) TestXorImmediate() {
	suite.operator.Bit_xor_immediate(0, 5, 7)
	suite.AssertRegisterEquals(0, 2)
	suite.AssertRegisterEquals(5, 5)
}

func (suite *OperatorSuite) TestNotImmediate() {
	suite.operator.Bit_not_immediate(0, 1)
	suite.AssertRegisterEquals(0, math.MaxUint32-1)
}

func (suite *OperatorSuite) TestLeftShiftImmediate() {
	suite.operator.Left_shift_immediate(0, 1, 3)
	suite.AssertRegisterEquals(0, 8)
	suite.AssertRegisterEquals(1, 1)

}

func (suite *OperatorSuite) TestRightShiftImmediate() {
	suite.operator.Right_shift_immediate(0, 5, 1, false)
	suite.AssertRegisterEquals(0, 2)
	suite.AssertRegisterEquals(5, 5)

	suite.memory.val = math.MaxUint32 - 2
	suite.operator.Load(5, 3, &suite.memory)
	suite.operator.Right_shift_immediate(0, 5, 1, true)
	suite.AssertRegisterEquals(0, math.MaxUint32-1)

}

func (suite *OperatorSuite) TestLoadWord() {
	suite.memory.val = math.MaxUint32
	suite.operator.Load_word(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint32)
}

func (suite *OperatorSuite) TestLoadHalfWord() {
	suite.memory.val = math.MaxUint16
	suite.operator.Load_halfword(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint32)
}

func (suite *OperatorSuite) TestLoadHalfWordUnsigned() {
	suite.memory.val = math.MaxUint16
	suite.operator.Load_halfword_unsigned(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint16)
}

func (suite *OperatorSuite) TestLoadByte() {
	suite.memory.val = math.MaxUint8
	suite.operator.Load_byte(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint32)
}

func (suite *OperatorSuite) TestLoadByteUnsigned() {
	suite.memory.val = math.MaxUint8
	suite.operator.Load_byte_unsigned(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint8)
}

func (suite *OperatorSuite) TestStoreWord() {
	assert := assert.New(suite.T())
	suite.memory.val = 0
	suite.operator.Store_word(10, 3, &suite.memory)

	assert.Equal(uint32(10), suite.memory.val)
}

func (suite *OperatorSuite) TestStoreHalfWord() {
	assert := assert.New(suite.T())
	suite.memory.val = 0
	suite.operator.Store_halfword(10, 3, &suite.memory)

	assert.Equal(uint32(10), suite.memory.val)
}

func (suite *OperatorSuite) TestStoreByte() {
	assert := assert.New(suite.T())
	suite.memory.val = 0
	suite.operator.Store_byte(10, 3, &suite.memory)

	assert.Equal(uint32(10), suite.memory.val)
}
