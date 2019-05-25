package adapters

import (
	"math"

	Memory "github.com/chenhowa/computer/lib/memory"
)

/*BasicMemory16Adapted adapts Memory.BasicMemory to use 32-bit addresses instead of 16-bit addresses, while
retaining (nearly) the same interface. For many of its methods, if an address is passed to BasicMemory16Adapted that is valid 32-bit but
not valid 16-bit, the method will panic.
*/
type BasicMemory16Adapted struct {
	memory *Memory.BasicMemory
}

/*Get returns the value stored in memory at `address`. Panics if
`address` is greater than math.MaxUint16
*/
func (m *BasicMemory16Adapted) Get(address uint32) uint32 {
	if address > math.MaxUint16 {
		panic("BasicMemory16Adapted: Get received too large of an address")
	}

	return m.memory.Get(uint16(address))
}

/*Set writes the lowest `bitsToWrite` bits of `val` to memory at `address`. Panics if
`address` is greater than math.MaxUint16
*/
func (m *BasicMemory16Adapted) Set(address uint32, val uint32, bitsToWrite uint) Memory.NumberOfBitsWritten {
	if address > math.MaxUint16 {
		panic("BasicMemory16Adapted: Set received too large of an address")
	}

	return m.memory.Set(uint16(address), val, bitsToWrite)
}

/*GetAddressSpaceSize returns the size of the address space that
this memory is valid for.
*/
func (m *BasicMemory16Adapted) GetAddressSpaceSize() uint {
	return m.memory.GetAddressSpaceSize()
}

/*MakeBasicMemory16Adapted is a constructor for BasicMemory16Adapted
 */
func MakeBasicMemory16Adapted(maxAddress uint16) BasicMemory16Adapted {
	basicMemory := Memory.MakeBasicMemory(maxAddress)
	memory := BasicMemory16Adapted{
		memory: &basicMemory,
	}
	return memory
}
