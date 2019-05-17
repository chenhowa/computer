package envManagers

/*NoOpExecManager presents an interface for managing
an outer execution environment, but that interface actually does
nothing*/
type NoOpExecManager struct {
}

/*ExecuteCall supposedly signals the outer executing environment
to execute a system call based on the values in the CPU registers,
but in actuality this does nothing.
*/
func (m *NoOpExecManager) ExecuteCall() {
	return
}
