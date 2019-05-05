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
	operator := df
	suite.executor = RiscVInstructionExecutor{
		operator: &operator,
	}
}
func (suite *InstructionExecutorSuite) TestAddImmediate() {
	assert := assert.New(suite.T())

}
