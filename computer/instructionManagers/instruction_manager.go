package instructionmanagers

/*PCInstructionManager supports 16 bit address space, but each location has 32 bits
 */
type PCInstructionManager struct {
	instructionAddress uint16
}

type managerReadMemory interface {
	Get(address uint16) uint32
}

/*GetCurrentInstructionAddress returns the address where the current instruction is stored
 */
func (manager *PCInstructionManager) GetCurrentInstructionAddress() uint16 {
	return manager.instructionAddress
}

/*GetCurrentInstruction gets the instruction stored at the current instruction address*/
func (manager *PCInstructionManager) GetCurrentInstruction(memory managerReadMemory) uint32 {
	return memory.Get(manager.GetCurrentInstructionAddress())
}

/*GetNextInstructionAddress gets the address of the instruction that is immediately AFTER the current instruction
 */
func (manager *PCInstructionManager) GetNextInstructionAddress(memory managerReadMemory) uint16 {
	// since each instruction is 4 bytes wide
	return manager.GetCurrentInstructionAddress() + 4
}

/*IncrementInstructionAddress updates the manager to point at the NextInstructionAddress*/
func (manager *PCInstructionManager) IncrementInstructionAddress() {
	manager.instructionAddress += 4 // notice that we increment by 4! At its base, the address space is byte-addressable, not bit-addressable.
}

/*AddOffsetForNextAddress updates the manager so the next Instruction Address is
<current instruction address> + `offset`
*/
func (manager *PCInstructionManager) AddOffsetForNextAddress(offset uint16) {
	manager.instructionAddress += (offset - 4)
}

/*LoadInstructionAddressForNextAddress updates the manager so that the next Instruction Address
is `newAddress`*/
func (manager *PCInstructionManager) LoadInstructionAddressForNextAddress(newAddress uint16) {
	manager.instructionAddress = (newAddress - 4)
}
