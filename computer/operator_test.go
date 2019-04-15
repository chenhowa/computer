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
	suite.operator = Operator{}
	suite.memory = MemoryMock{}
}

func TestRun(t *testing.T) {
	suite.Run(t, new(OperatorSuite))
}

func (suite *OperatorSuite) TestLoadStore() {
	var assert = assert.New(suite.T())

	var zero uint32 = 0
	assert.Equal(suite.operator.registers[0], zero)

	var start_value uint32 = 10
	suite.memory.val = start_value
	assert.Equal(suite.memory.val, start_value)

	suite.operator.load(0, 5, &suite.memory)
	assert.Equal(suite.operator.registers[0], start_value)

	suite.operator.load(15, 5, &suite.memory)
	assert.Equal(suite.operator.registers[15], start_value)

	suite.operator.store(1, 5, &suite.memory)
	assert.Equal(suite.memory.val, zero)
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

}
