package execution

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OperatorSuite struct {
	suite.Suite
	operator adaptedOperator
	memory   MemoryMock
}

type MemoryMock struct {
	val uint32
}

func (m *MemoryMock) Get(address uint32) uint32 {
	return m.val
}

func (m *MemoryMock) Set(address uint32, val uint32, bits uint) {
	m.val = val
}

func (suite *OperatorSuite) SetupTest() {
	suite.operator = makeAdaptedOperator([16]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})

	suite.memory = MemoryMock{}
}

func (suite *OperatorSuite) AssertRegisterEquals(reg uint16, val uint32) {
	assert := assert.New(suite.T())
	actual := suite.operator.get(uint(reg))
	assert.Equal(val, actual)
}

func TestRun(t *testing.T) {
	suite.Run(t, new(OperatorSuite))
}

func (suite *OperatorSuite) TestAdd() {
	var assert = assert.New(suite.T())

	var startValue uint32 = 10
	suite.memory.val = startValue
	suite.operator.loadWord(0, 5, &suite.memory)
	suite.operator.loadWord(15, 5, &suite.memory)

	suite.operator.add(5, 0, 15)
	assert.Equal(suite.operator.get(5), startValue*2)
}

func (suite *OperatorSuite) TestSub() {
	var assert = assert.New(suite.T())

	suite.operator.sub(0, 3, 1)
	suite.operator.storeWord(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(2))
}

func (suite *OperatorSuite) TestBitAnd() {
	var assert = assert.New(suite.T())

	suite.operator.and(0, 1, 3)
	suite.operator.storeWord(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(1))
}

func (suite *OperatorSuite) TestBitOr() {
	var assert = assert.New(suite.T())

	suite.operator.or(0, 1, 2)
	suite.operator.storeWord(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(3))
}

func (suite *OperatorSuite) TestBitXor() {
	suite.operator.xor(0, 5, 3) // 0b101 ^ 0b011 = 6
	suite.AssertRegisterEquals(0, 6)
	suite.AssertRegisterEquals(3, 3)
	suite.AssertRegisterEquals(5, 5)
}

func (suite *OperatorSuite) TestMultiply() {
	var assert = assert.New(suite.T())
	suite.operator.multiply(0, 2, 3)
	suite.operator.storeWord(0, 5, &suite.memory)
	assert.Equal(suite.memory.val, uint32(6))
}

func (suite *OperatorSuite) TestDivide() {
	suite.operator.divide(0, 1, 5, 3)
	suite.AssertRegisterEquals(0, uint32(1))
	suite.AssertRegisterEquals(1, uint32(2))
}

func (suite *OperatorSuite) TestAddImmediate() {
	suite.operator.addImmediate(0, 1, 10)
	suite.AssertRegisterEquals(0, 11)
	suite.AssertRegisterEquals(1, 1)
}

func (suite *OperatorSuite) TestAndImmediate() {
	suite.operator.andImmediate(0, 2, 7)
	suite.AssertRegisterEquals(0, 2)
	suite.AssertRegisterEquals(2, 2)
}

func (suite *OperatorSuite) TestOrImmediate() {
	suite.operator.orImmediate(0, 5, 2)
	suite.AssertRegisterEquals(0, 7)
	suite.AssertRegisterEquals(5, 5)
}

func (suite *OperatorSuite) TestXorImmediate() {
	suite.operator.xorImmediate(0, 5, 7)
	suite.AssertRegisterEquals(0, 2)
	suite.AssertRegisterEquals(5, 5)
}

func (suite *OperatorSuite) TestLeftShiftImmediate() {
	suite.operator.leftShiftImmediate(0, 1, 3)
	suite.AssertRegisterEquals(0, 8)
	suite.AssertRegisterEquals(1, 1)

}

func (suite *OperatorSuite) TestRightShiftImmediate() {
	suite.operator.rightShiftImmediate(0, 5, 1, false) // THIS IS CORRECT. THIS IS JUST AN ADAPTER FOR OPERATOR. IT DOESN"T DO ANYTHING RISC WISE
	suite.AssertRegisterEquals(0, 2)
	suite.AssertRegisterEquals(5, 5)

	suite.memory.val = math.MaxUint32 - 2
	suite.operator.loadWord(5, 3, &suite.memory)
	suite.operator.rightShiftImmediate(0, 5, 1, true) // THIS IS CORRECT. THIS IS JUST AN ADAPTER FOR OPERATOR. IT DOESN"T DO ANYTHING RISC WISE
	suite.AssertRegisterEquals(0, math.MaxUint32-1)

}

func (suite *OperatorSuite) TestLoadWord() {
	suite.memory.val = math.MaxUint32
	suite.operator.loadWord(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint32)
}

func (suite *OperatorSuite) TestLoadHalfWord() {
	suite.memory.val = math.MaxUint16
	suite.operator.loadHalfWord(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint32)
}

func (suite *OperatorSuite) TestLoadHalfWordUnsigned() {
	suite.memory.val = math.MaxUint16
	suite.operator.loadHalfWordUnsigned(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint16)
}

func (suite *OperatorSuite) TestLoadByte() {
	suite.memory.val = math.MaxUint8
	suite.operator.loadByte(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint32)
}

func (suite *OperatorSuite) TestLoadByteUnsigned() {
	suite.memory.val = math.MaxUint8
	suite.operator.loadByteUnsigned(0, 3, &suite.memory)
	suite.AssertRegisterEquals(0, math.MaxUint8)
}

func (suite *OperatorSuite) TestStoreWord() {
	assert := assert.New(suite.T())
	suite.memory.val = 0
	suite.operator.storeWord(10, 3, &suite.memory)

	assert.Equal(uint32(10), suite.memory.val)
}

func (suite *OperatorSuite) TestStoreHalfWord() {
	assert := assert.New(suite.T())
	suite.memory.val = 0
	suite.operator.storeHalfWord(10, 3, &suite.memory)

	assert.Equal(uint32(10), suite.memory.val)
}

func (suite *OperatorSuite) TestStoreByte() {
	assert := assert.New(suite.T())
	suite.memory.val = 0
	suite.operator.storeByte(10, 3, &suite.memory)

	assert.Equal(uint32(10), suite.memory.val)
}
