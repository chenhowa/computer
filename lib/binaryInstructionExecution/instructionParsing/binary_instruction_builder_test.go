package instructionParsing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BuilderSuite struct {
	suite.Suite
	builder BinaryInstructionBuilder
}

func TestBuilderSuite(t *testing.T) {
	suite.Run(t, new(BuilderSuite))
}

func (suite *BuilderSuite) SetupTest() {
	suite.builder = MakeInstructionBuilder(32)
}

func (suite *BuilderSuite) TestEmptyBuild() {
	assert := assert.New(suite.T())

	assert.Equal(uint(0), suite.builder.Build())
}

func (suite *BuilderSuite) TestPartialFillBuild() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(2, 0) // 00
	suite.builder.AddNextXBits(2, 3) // 11

	assert.Equal(uint(12), suite.builder.Build())

	suite.builder.AddNextXBits(1, 0) // 0
	suite.builder.AddNextXBits(3, 7) // 111

	assert.Equal(uint(236), suite.builder.Build())
}

func (suite *BuilderSuite) TestOverflowBuild() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(30, 0)
	suite.builder.AddNextXBits(3, 7)

	assert.Equal(uint(3221225472), suite.builder.Build())

	suite.builder.AddNextXBits(3, 7)
	assert.Equal(uint(3221225472), suite.builder.Build())
}

func (suite *BuilderSuite) TestTakesLowestXBits() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(2, 15)

	assert.Equal(uint(3), suite.builder.Build())
}
