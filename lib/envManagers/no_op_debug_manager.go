package envManagers

/*NoOpDebugManager presents an interface for managing
an outer debugging environment, but that interface actually does
nothing*/
type NoOpDebugManager struct {
}

/*DebugBreak supposedly signals the outer debugging environment
to pause execution at a breakpoint, but in actuality it does nothing
*/
func (m *NoOpDebugManager) DebugBreak() {
	return
}
