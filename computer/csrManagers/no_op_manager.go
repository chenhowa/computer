package csrManagers

/*NoOpManager is an interface to reading/writing the control status registers (CSRs), where there is only
one register, that does not secretly interact with anything else.*/
type NoOpManager struct {
	register uint32
}

/*Get returns the value stored in the single register of NoOpManager, regardless
of the value of `register`
*/
func (m *NoOpManager) Get(register uint) uint32 {
	return m.register
}

/*Set writes the value `val` to the single register of NoOpRegister, regardless
of the value of `register`*/
func (m *NoOpManager) Set(register uint, val uint32) {
	m.register = val
}
