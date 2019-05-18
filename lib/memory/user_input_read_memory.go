package memory

/*UserInputReadMemory adapts memory so that when the caller tries to read from an address,
the value at the address is obtained externally, through the command line or the network,
and returned to the caller, and the value is then recorded in memory. If necessary,
the memory will block while waiting for the external value to be sent.
*/
type UserInputReadMemory struct {
	memory memory
	source source
}

type memory interface {
	set(address uint16, value uint32)
	get(address uint16) uint32
}

type source interface {
	get(address uint16, existingVal uint32) uint32
}

/*Get will attempt to get a value from an external `source`, while passing it the existing value at
the address, in case the external wants to choose the existing value. If necessary, Get will
block as it waits for the value. Once the value is obtained,
it will write the value to memory at the 16-bit `address`, and return the value
to the caller, so that we can pretend the value was always there in memory.
*/
func (m *UserInputReadMemory) Get(address uint16) uint32 {
	val := m.source.get(address, m.memory.get(address))
	m.memory.set(address, val)
	return val
}
