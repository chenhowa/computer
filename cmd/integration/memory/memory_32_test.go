package memory

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type Memory32Suite struct {
	suite.Suite
	memory *Memory32
	sink   *MockErrorHandler
}

func TestMemory32Suite(t *testing.T) {
	suite.Run(t, new(Memory32Suite))
}

func (suite *Memory32Suite) SetupTest() {
	sink := MockErrorHandler{}
	memory := MakeMemory32(20, &sink)
	suite.memory = &memory
	suite.sink = &sink
}

func (suite *Memory32Suite) AssertMemoryAtAddressIs(address uint32, val uint32) {
	assert.Equal(suite.T(), val, suite.memory.Get(address))
}

func (suite *Memory32Suite) AssertNoMemoryErrors() {
	suite.sink.AssertNotCalled(suite.T(), "Handle", mock.Anything)
}

type MockErrorHandler struct {
	mock.Mock
}

func (h *MockErrorHandler) Handle(message interface{}) {
	h.Called(message)
}

func (suite *Memory32Suite) TestGet() {
	suite.AssertMemoryAtAddressIs(0, 0)
	suite.memory.Set(0, math.MaxUint16, 16)
	suite.AssertNoMemoryErrors()
	suite.AssertMemoryAtAddressIs(0, math.MaxUint16)
	suite.AssertMemoryAtAddressIs(1, math.MaxUint8)
}

func (suite *Memory32Suite) TestSet() {
	suite.AssertMemoryAtAddressIs(20, 0)
	suite.memory.Set(20, math.MaxUint16, 16)
	suite.AssertNoMemoryErrors()
	suite.AssertMemoryAtAddressIs(20, math.MaxUint8)
}

func (suite *Memory32Suite) TestSet_Fail() {
	suite.AssertMemoryAtAddressIs(21, 0)
	suite.memory.Set(21, math.MaxUint16, 16)
	suite.AssertMemoryAtAddressIs(21, math.MaxUint8)
}
