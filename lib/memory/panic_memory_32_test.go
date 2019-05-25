package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PanicMemorySuite struct {
	suite.Suite
	sink   *MockErrorSink
	memory *PanicMemory32
}

func TestPanicMemorySuite(t *testing.T) {
	suite.Run(t, new(PanicMemorySuite))
}

func (suite *PanicMemorySuite) SetupTest() {
	sink := MockErrorSink{}
	suite.sink = &sink

	basicMemory := MockBasicMemory{}

	memory := PanicMemory32{}
	memory.memory = &basicMemory
	memory.errorSink = suite.sink
	suite.memory = &memory
}

type MockErrorSink struct {
	mock.Mock
}

func (es *MockErrorSink) handle(message interface{}) {
	es.Called(message)
}

type MockBasicMemory struct {
}

func (m *MockBasicMemory) Get(address uint32) uint32 {
	if address > 10 {
		panic(1)
	} else {
		return 5
	}
}

func (m *MockBasicMemory) Set(address uint32, val uint32, bitsToWrite uint) NumberOfBitsWritten {
	if address > 10 {
		panic(1)
	} else {
		return 20
	}
}

func (m *MockBasicMemory) GetAddressSpaceSize() uint {
	return 11
}

func (suite *PanicMemorySuite) TestSet() {
	assert := assert.New(suite.T())
	result := suite.memory.Set(5, 5, 5)
	assert.Equal(NumberOfBitsWritten(20), result)
	suite.sink.AssertNotCalled(suite.T(), "handle", mock.Anything)

	suite.sink.On("handle", mock.Anything)
	result = suite.memory.Set(11, 5, 5)
	assert.Equal(NumberOfBitsWritten(0), result)
	suite.sink.AssertCalled(suite.T(), "handle", mock.Anything)
}

func (suite *PanicMemorySuite) TestGet() {
	assert := assert.New(suite.T())
	result := suite.memory.Get(5)
	assert.Equal(uint32(5), result)
	suite.sink.AssertNotCalled(suite.T(), "handle", mock.Anything)

	suite.sink.On("handle", mock.Anything)
	result = suite.memory.Get(11)
	assert.Equal(uint32(0), result)
	suite.sink.AssertCalled(suite.T(), "handle", mock.Anything)
}
