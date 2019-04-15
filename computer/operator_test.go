package computer

import (
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

func (m *MemoryMock) Set(address uint16, val uint32) {
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

	suite.operator.load(0, 5, &suite.memory)
	assert.Equal(suite.operator.registers[0], start_value)

	suite.operator.load(15, 5, &suite.memory)
	assert.Equal(suite.operator.registers[15], start_value)

	suite.operator.store(1, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(1))
}

func (suite *OperatorSuite) TestAdd() {
	var assert = assert.New(suite.T())

	var start_value uint32 = 10
	suite.memory.val = start_value
	suite.operator.load(0, 5, &suite.memory)
	suite.operator.load(15, 5, &suite.memory)

	suite.operator.add(5, 0, 15)
	assert.Equal(suite.operator.registers[5], start_value*2)
}

func (suite *OperatorSuite) TestSub() {
	var assert = assert.New(suite.T())

	suite.operator.sub(0, 3, 1)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(2))
}

func (suite *OperatorSuite) TestBitAnd() {
	var assert = assert.New(suite.T())

	suite.operator.bit_and(0, 1, 3)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(1))
}

func (suite *OperatorSuite) TestBitOr() {
	var assert = assert.New(suite.T())

	suite.operator.bit_or(0, 1, 2)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(3))
}

func (suite *OperatorSuite) TestLeftShift() {
	var assert = assert.New(suite.T())

	suite.operator.left_shift(0, 1)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(2))
}

func (suite *OperatorSuite) TestRightShiftLogical() {
	var assert = assert.New(suite.T())

	suite.operator.right_shift(0, 2, false)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(1))
}

func (suite *OperatorSuite) TestRightShiftArithmetic() {
	var assert = assert.New(suite.T())
	suite.memory.val = 1 >> 31
	suite.operator.load(1, 5, &suite.memory)
	suite.operator.right_shift(0, 1, true)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32((1>>30)|(1>>31)))
}

func (suite *OperatorSuite) TestMultiply() {
	var assert = assert.New(suite.T())
	suite.operator.multiply(0, 2, 3)
	suite.operator.store(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(6))
}

func (suite *OperatorSuite) TestDivide() {
	suite.operator.divide(0, 1, 5, 3)
	suite.AssertRegisterEquals(0, uint32(1))
	suite.AssertRegisterEquals(1, uint32(2))
}
