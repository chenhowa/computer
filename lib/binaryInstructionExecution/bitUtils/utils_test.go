package bitUtils

import (
	"math"
	"testing"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsSuite struct {
	suite.Suite
}

func (suite *UtilsSuite) SetupTest() {

}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}

func (suite *UtilsSuite) TestSignExtension() {
	assert := assert.New(suite.T())
	//Test sign extension with 1
	assert.Equal(uint32(math.MaxUint32), SignExtendUint32WithBit(1, 0))
	assert.Equal(uint32(math.MaxUint32-1), SignExtendUint32WithBit(2, 1))
	assert.Equal(uint32(math.MaxUint32-2), SignExtendUint32WithBit(5, 2))
	assert.Equal(uint32(math.MaxUint32), SignExtendUint32WithBit(5, 0))

	// Test sign extension with 0
	assert.Equal(uint32(1), SignExtendUint32WithBit(1, 1))
	assert.Equal(uint32(0), SignExtendUint32WithBit(2, 0))
	assert.Equal(uint32(2), SignExtendUint32WithBit(2, 2))
	assert.Equal(uint32(1), SignExtendUint32WithBit(5, 1))
	assert.Equal(uint32(5), SignExtendUint32WithBit(5, 3))
}

func (suite *UtilsSuite) TestSignExtension2() {
	assert := assert.New(suite.T())

	assert.Equal(uint32((1<<11)-1), SignExtendUint32WithBit(KeepBitsInInclusiveRange(math.MaxUint32, 0, 10), 11))
}

func (suite *UtilsSuite) TestKeepBitsInclusiveRange() {
	assert := assert.New(suite.T())

	assert.Equal(uint32(4), KeepBitsInInclusiveRange(7, 2, 2))
	assert.Equal(uint32(6), KeepBitsInInclusiveRange(7, 1, 2))
	assert.Equal(uint32(3), KeepBitsInInclusiveRange(7, 0, 1))

}
