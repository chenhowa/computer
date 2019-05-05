package execution

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InstructionExecutorSuite struct {
	suite.Suite
	executor RiscVInstructionExecutor
}

func TestInstructionExecutorSuite(t *testing.T) {
	suite.Run(t, new(InstructionExecutorSuite))
}

func (suite *InstructionExecutorSuite) SetupTest() {
	operator := Operator{
		registers: [16]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	}
	suite.executor = RiscVInstructionExecutor{
		operator: &operator,
	}
}

func (suite *InstructionExecutorSuite) TestSignExtension() {
	assert := assert.New(suite.T())

	//Test sign extension with 1
	assert.Equal(uint32(math.MaxUint32), signExtendUint32WithBit(1, 0))
	assert.Equal(uint32(math.MaxUint32-1), signExtendUint32WithBit(2, 1))
	assert.Equal(uint32(math.MaxUint32-2), signExtendUint32WithBit(5, 2))
	assert.Equal(uint32(math.MaxUint32), signExtendUint32WithBit(5, 0))

	// Test sign extension with 0
	assert.Equal(uint32(1), signExtendUint32WithBit(1, 1))
	assert.Equal(uint32(0), signExtendUint32WithBit(2, 0))
	assert.Equal(uint32(2), signExtendUint32WithBit(2, 2))
	assert.Equal(uint32(1), signExtendUint32WithBit(5, 1))
	assert.Equal(uint32(5), signExtendUint32WithBit(5, 3))
}

func (suite *InstructionExecutorSuite) TestAddImmediate() {
	assert := assert.New(suite.T())

}
