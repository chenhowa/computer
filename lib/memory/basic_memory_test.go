package memory

import (
	"math"
	"testing"

	Utils "github.com/chenhowa/computer/lib/binaryInstructionExecution/bitUtils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BasicMemorySuite struct {
	suite.Suite
}

func TestBasicMemorySuite(t *testing.T) {
	suite.Run(t, new(BasicMemorySuite))
}

func (suite *BasicMemorySuite) SetupTest() {

}

func (suite *BasicMemorySuite) TestSet_Basic_1() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(0, math.MaxUint32, 8)
	assert.Equal(uint32(math.MaxUint8), memory.Get(0))
	assert.Equal(uint32(0), memory.Get(1))
}

func (suite *BasicMemorySuite) TestSet_Basic_2() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(0, math.MaxUint32, 32)
	assert.Equal(uint32(math.MaxUint32), memory.Get(0))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 23), memory.Get(1))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 15), memory.Get(2))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 7), memory.Get(3))
	assert.Equal(uint32(0), memory.Get(4))
}

func (suite *BasicMemorySuite) TestSet_Basic_3() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(0, math.MaxUint32, 10)
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 9), memory.Get(0))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 1), memory.Get(1))
	assert.Equal(uint32(0), memory.Get(2))
}

func (suite *BasicMemorySuite) TestSet_Basic_4() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(4, uint32(math.MaxUint8-1), 2)
	assert.Equal(uint32(2), memory.Get(4))

	memory.Set(0, math.MaxUint32, 33)
	assert.Equal(uint32(math.MaxUint32), memory.Get(0))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 23)+Utils.KeepBitsInInclusiveRange(math.MaxUint32, 25, 25), memory.Get(1))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 15)+Utils.KeepBitsInInclusiveRange(math.MaxUint32, 17, 17), memory.Get(2))
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 0, 7)+Utils.KeepBitsInInclusiveRange(math.MaxUint32, 9, 9), memory.Get(3))
	assert.Equal(uint32(2), memory.Get(4))
}

func (suite *BasicMemorySuite) TestSet_Near_Max_Address_1() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(19, Utils.KeepBitsInInclusiveRange(math.MaxUint32, 4, 31), 16)
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 4, 7), memory.Get(19))
}

func (suite *BasicMemorySuite) TestSet_Near_Max_Address_2() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(18, Utils.KeepBitsInInclusiveRange(math.MaxUint32, 4, 31), 16)
	assert.Equal(Utils.KeepBitsInInclusiveRange(math.MaxUint32, 4, 15), memory.Get(18))
}

func (suite *BasicMemorySuite) TestSet_Panic() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	assert.Panics(func() {
		memory.Set(20, 100, 5)
	})
	assert.Panics(func() {
		memory.Set(21, 100, 5)
	})
}

func (suite *BasicMemorySuite) TestGet() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	memory.Set(4, uint32(math.MaxUint8-1), 2)
	assert.Equal(uint32(2), memory.Get(4))
}

func (suite *BasicMemorySuite) TestGet_Panic() {
	memory := MakeBasicMemory(19)
	assert := assert.New(suite.T())

	assert.Equal(uint(20), memory.GetAddressSpaceSize())
	assert.Equal(uint32(0), memory.Get(0))
	assert.Equal(uint32(0), memory.Get(19))
	assert.Panics(func() {
		memory.Get(20)
	})
	assert.Panics(func() {
		memory.Get(21)
	})

}
