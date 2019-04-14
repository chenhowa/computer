package computer


type RAM struct {
	memory [65536]uint32
}

func (r RAM) Get(address uint16) uint32 {
	return r.memory[address]
}

func (r RAM) Set(address uint16, value uint32) {
	r.memory[address] = value
}




type ROM struct {
	memory [65536]uint32
}

func (r ROM) Get(address uint16) uint32 {
	return r.memory[address]
}

