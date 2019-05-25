package memory

/*PanicMemory32 assumes that the underlying memory may panic, and as a side effect,
it handles it by informing an error sink of the error*/
type PanicMemory32 struct {
	memory    basicMemory
	errorSink errorSink
}

type errorSink interface {
	Handle(message interface{})
}

/*MakePanicMemory32 is a constructor for PanicMemory32*/
func MakePanicMemory32(memory basicMemory, sink errorSink) PanicMemory32 {
	panicMemory := PanicMemory32{
		memory:    memory,
		errorSink: sink,
	}

	return panicMemory
}

/*Get attempts to get the value from memory at `address`
If the attempt panics, the errorSink will be informed*/
func (m *PanicMemory32) Get(address uint32) (v uint32) {
	defer func() {
		if r := recover(); r != nil {
			v = 0
			m.errorSink.Handle(r)
		}
	}()

	val := m.memory.Get(address)

	return val
}

/*Set attempts to write `bitsToWrite` bits of `val` to memory at `address`
If the attempt panics, the errorSink will be informed*/
func (m *PanicMemory32) Set(address uint32, val uint32, bitsToWrite uint) (n NumberOfBitsWritten) {
	defer func() {
		if r := recover(); r != nil {
			n = 0
			m.errorSink.Handle(r)
		}
	}()

	numBits := m.memory.Set(address, val, bitsToWrite)

	return numBits
}
