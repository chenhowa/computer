package sources

/*MemorySource is a source that simply accepts that the value in memory
is the correct value
*/
type MemorySource struct {
}

/*Get simply returns the existing memory value as the value in memory*/
func (s *MemorySource) Get(address uint16, existingVal uint32) uint32 {
	return existingVal
}
