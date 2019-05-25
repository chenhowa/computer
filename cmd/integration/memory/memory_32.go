package memory

import (
	"math"

	Adapters "github.com/chenhowa/computer/cmd/integration/memory/adapters"
	Memory "github.com/chenhowa/computer/lib/memory"
)

/*Memory32 represents 32-bit memory that can be written to and read from.
Internally, however, AT MOST only the first 16-bits of memory are valid. This is because
the intent is to simulate 16-bit memory, even though the RISC-V 32-bit instruction
set is being used. After all, who can guarantee the amount of RAM available to the CPU?
*/
type Memory32 struct {
	sink   errorSink
	memory *Memory.PanicMemory32
}

type errorSink interface {
	Handle(message interface{})
}

/*MakeMemory32 is a constructor for Memory32. It accepts an errorSink that
will be responsible for handling any errors that occur while fetching or writing memory*/
func MakeMemory32(maxAddress uint16, sink errorSink) Memory32 {
	memory := Memory32{}
	errorSink := sink

	basicMemory := Adapters.MakeBasicMemory16Adapted(math.MaxUint16)

	panicMemory := Memory.MakePanicMemory32(&basicMemory, errorSink)
	memory.memory = &panicMemory

	return memory
}

/*Get fetches a 32-bit value from the `address` in memory*/
func (m *Memory32) Get(address uint32) uint32 {
	return m.memory.Get(address)
}

/*Set writes the first `numBits` of a 32-bit `val` to memory at `address`,
and returns the number of bits successfully written*/
func (m *Memory32) Set(address uint32, val uint32, numBits uint) Memory.NumberOfBitsWritten {
	result := m.memory.Set(address, val, numBits)
	return result
}
