package instructionmanagers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MemoryMock struct {
	mock.Mock
}

func (m *MemoryMock) Get(address uint16) uint32 {
	m.Called(address)
	return 2000
}

type InstructionManagerSuite struct {
	suite.Suite
	manager *PCInstructionManager
	memory  *MemoryMock
}

func TestBuilderSuite(t *testing.T) {
	suite.Run(t, new(InstructionManagerSuite))
}

func (suite *InstructionManagerSuite) SetupTest() {
	manager := MakePCInstructionManager(15)
	suite.manager = &manager

	memory := MemoryMock{}
	suite.memory = &memory
}

func (suite *InstructionManagerSuite) TestGetCurrentInstructionAddress() {
	assert := assert.New(suite.T())

	assert.Equal(uint16(11), suite.manager.GetCurrentInstructionAddress())
}

func (suite *InstructionManagerSuite) TestGetCurrentInstruction() {
	assert := assert.New(suite.T())
	suite.memory.On("Get", uint16(11))
	assert.Equal(uint32(2000), suite.manager.GetCurrentInstruction(suite.memory))
	suite.memory.AssertCalled(suite.T(), "Get", uint16(11))
}

func (suite *InstructionManagerSuite) TestGetNextInstructionAddress() {
	assert := assert.New(suite.T())

	assert.Equal(uint16(15), suite.manager.GetNextInstructionAddress())
}

func (suite *InstructionManagerSuite) TestIncrementInstructionAddress() {
	assert := assert.New(suite.T())

	suite.manager.IncrementInstructionAddress()
	assert.Equal(uint16(15), suite.manager.GetCurrentInstructionAddress())

}

func (suite *InstructionManagerSuite) TestAddOffsetForNextAddress() {
	assert := assert.New(suite.T())

	suite.manager.AddOffsetForNextAddress(20)
	assert.Equal(uint16(11+20-4), suite.manager.GetCurrentInstructionAddress())
	assert.Equal(uint16(11+20), suite.manager.GetNextInstructionAddress())
	suite.manager.IncrementInstructionAddress()
	assert.Equal(uint16(11+20), suite.manager.GetCurrentInstructionAddress())
}

func (suite *InstructionManagerSuite) TestLoadInstructionAddressForNextAddress() {
	assert := assert.New(suite.T())
	suite.manager.LoadInstructionAddressForNextAddress(11 + 20)
	assert.Equal(uint16(11+20-4), suite.manager.GetCurrentInstructionAddress())
	assert.Equal(uint16(11+20), suite.manager.GetNextInstructionAddress())
	suite.manager.IncrementInstructionAddress()
	assert.Equal(uint16(11+20), suite.manager.GetCurrentInstructionAddress())

}
