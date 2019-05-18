package memory

import "errors"

/*CompositeMemory32 presents methods for reading and writing from a 32-bit address
memory space. It does not enforce the idea that each consecutive address addresses
consecutive bytes. This struct supports reading up to a doubleword from memory, and no more.
*/
type CompositeMemory32 struct {
	memory        basicMemory
	furtherMemory *CompositeMemory32
}

/*NumberOfBitsWritten represents the number of bits written by a call to Set*/
type NumberOfBitsWritten uint

type basicMemory interface {
	Get(address uint32) uint32
	Set(address uint32, val uint32, bitsToWrite uint) NumberOfBitsWritten
	GetAddressSpaceSize() uint
}

/*GetAddressSpaceSize returns the first address that this memory
does not support; that is, it returns (MaxAddress + 1). Note that
CompositeMemory32 will always start its address support from 0*/
func (m *CompositeMemory32) GetAddressSpaceSize() uint {
	return m.memory.GetAddressSpaceSize() + m.furtherMemory.GetAddressSpaceSize()
}

/*Get attempts to retrieve the doubleword value from `address` in memory.
However, if `address` is outside the range of addressable memory,
th call to Get will return an error, and it is up to the caller to handle what happens.
*/
func (m *CompositeMemory32) Get(address uint32) (uint32, error) {
	currentMemorySize := uint32(m.memory.GetAddressSpaceSize())
	if !(address < uint32(m.GetAddressSpaceSize())) {
		return 0, errors.New("Get: Requested address is not addressable in memory")
	} else if address < currentMemorySize {
		return m.memory.Get(address), nil
	} else {
		return m.furtherMemory.Get(address - currentMemorySize)
	}
}

/*Set attempts to write the lowest `bitsToWrite` bits of `val` to memory, starting at `address`,
and starting with the LSB of `val`. Attempting to write more than 32 bits will still write only 32 bits. If writing
these bits to memory causes attempted writes to invalid addresses, Set will return an error
to the caller to handle. */
func (m *CompositeMemory32) Set(address uint32, val uint32, bitsToWrite uint) error {
	if bitsToWrite == 0 {
		return nil
	}

	currentMemorySize := uint32(m.memory.GetAddressSpaceSize())
	if !(address < uint32(m.GetAddressSpaceSize())) {
		return errors.New("Set: Requested address is not addressable in memory")
	} else if address < currentMemorySize {
		/*If the write did not complete within current memory, we will attempt to complete it
		in the `furtherMemory`, if possible*/
		bitsWritten := m.memory.Set(address, val, bitsToWrite)
		return m.furtherMemory.Set(0, val>>bitsWritten, bitsToWrite-uint(bitsWritten))
	} else {
		return m.furtherMemory.Set(address-currentMemorySize, val, bitsToWrite)
	}
}
