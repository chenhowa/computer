


package computer

import (
	"testing"
	"github.com/stretchr/testify/assert"
);

type MemoryMock struct {
	val uint32
}

func (m *MemoryMock) Get(address uint16) uint32 {
	return m.val
}

func (m *MemoryMock) Set(address uint16, val uint32) {
	m.val = val
} 

func TestLoadStore(t *testing.T) {
	var operator Operator =  Operator {};
	var assert = assert.New(t);

	var zero uint32 = 0;
	assert.Equal(operator.registers[0], zero);

	var start_value uint32 = 10
	memory := MemoryMock { val: start_value }
	assert.Equal(memory.val, start_value);

	operator.load(0, 5, &memory);
	assert.Equal(operator.registers[0], start_value);
	
	operator.load(15, 5, &memory);
	assert.Equal(operator.registers[15], start_value);

	operator.store(1, 5, &memory);
	assert.Equal(memory.val, zero);
}	

func TestAdd(t *testing.T) {
	var operator Operator =  Operator {};
	var assert = assert.New(t);

	var start_value uint32 = 10
	memory := MemoryMock { val: start_value }
	operator.load(0, 5, &memory);
	operator.load(15, 5, &memory);

	operator.add(5, 0, 15);
	assert.Equal(operator.registers[5], start_value * 2);
}