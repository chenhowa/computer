package execution

import (
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
func (suite *InstructionExecutorSuite) TestAddImmediate() {
	assert := assert.New(suite.T())

}
