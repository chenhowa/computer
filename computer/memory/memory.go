package memory

/*RAM IS WRITTEN INCORRECTLY. THE ADDRESSES ARE BYTE-ADDRESSABLE. SO A
16 BIT ADDRESS SPACE ADDRESSES 65536 BYTES, not 65536 WORDS!

*/
type RAM struct {
	memory [65536]uint32
}

func (r RAM) Get(address uint16) uint32 {
	return r.memory[address]
}

func (r RAM) Set(address uint16, value uint32) {
	r.memory[address] = value
}

/*ROM IS WRITTEN INCORRECTLY. THE ADDRESSES ARE BYTE-ADDRESSABLE. SO A
16 BIT ADDRESS SPACE ADDRESSES 65536 BYTES, not 65536 WORDS!

*/
type ROM struct {
	memory [65536]uint32
}

func (r ROM) Get(address uint16) uint32 {
	return r.memory[address]
}
