package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func TestUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}

func (suite *UtilSuite) SetupTest() {
}

func (suite *UtilSuite) TestIsNumericConstant_True() {
	assert := assert.New(suite.T())

	assert.Equal(true, isNumericConstant("0"))
	assert.Equal(true, isNumericConstant("1"))
	assert.Equal(true, isNumericConstant("1000"))
	assert.Equal(true, isNumericConstant("1,000"))
	assert.Equal(true, isNumericConstant("1111111"))
	assert.Equal(true, isNumericConstant("1,111,111"))
	assert.Equal(true, isNumericConstant("982838"))
}

func (suite *UtilSuite) TestIsNumericConstant_False() {
	assert := assert.New(suite.T())

	assert.Equal(false, isNumericConstant("01"))
	assert.Equal(false, isNumericConstant("0011"))
	assert.Equal(false, isNumericConstant("10,00"))
	assert.Equal(false, isNumericConstant("100,0"))
	assert.Equal(false, isNumericConstant("1000,"))
	assert.Equal(false, isNumericConstant(",1000"))
}
