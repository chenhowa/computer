package computer

/*PCInstructionManager supports 16 bit address space, but each location has 32 bits
 */
type PCInstructionManager struct {
	instructionAddress uint16
}

type managerReadMemory interface {
	Get(address uint16) uint32
}

func (manager *PCInstructionManager) getInstructionAddress() uint16 {
	return manager.instructionAddress
}

func (manager *PCInstructionManager) getInstruction(memory managerReadMemory) uint32 {
	return memory.Get(manager.instructionAddress)
}

func (manager *PCInstructionManager) incrementInstructionAddress() {
	manager.instructionAddress += 4 // notice that we increment by 4! At its base, the address space is byte-addressable.
}

func (manager *PCInstructionManager) addOffset(offset uint16) {
	manager.instructionAddress += offset
}

func (manager *PCInstructionManager) loadInstructionAddress(newAddress uint16) {
	manager.instructionAddress = newAddress
}
