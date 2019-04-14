


package computer

import "testing";

type MemoryMock struct {
	val uint32
}

func (m *MemoryMock) Get(address uint16) uint32 {
	return m.val
}

func (m *MemoryMock) Set(address uint16, val uint32) {
	m.val = val
} 

func TestAdd(t *testing.T) {
	var operator Operator =  Operator {};
	if operator.registers[0] != 0 {
		// Sanity check on initialization.
		t.Errorf("Expected 0, got %d", operator.registers[0] );
	}

	var start_value uint32 = 10
	memory := MemoryMock { val: start_value }
	if memory.val != start_value {
		t.Errorf("Expected %d, got %d", start_value, memory.val);
	}

	operator.load(0, 5, &memory);
	if operator.registers[0] != start_value {
		t.Errorf("Expected %d, got %d", start_value, operator.registers[0]);
	}

	operator.load(15, 5, &memory);
	if operator.registers[15] != start_value {
		t.Errorf("Expected %d, got %d", start_value, operator.registers[15])
	}

	operator.store(1, 5, &memory);
	if memory.val != 0 {
		t.Errorf("Expected %d, got %d", 0, memory.val);
	}
}