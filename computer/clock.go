package computer


/*
	This simulates the clock of a computer. Each instruction is executed, and the next one fetched,
	according to the timing of the clock. This is a useful part of the simulation because it allows us to
	speed up and slow down execution of the computer for observation and debugging.

	The clock itself can be "turned off" so that execution proceeds as quickly as possible.
*/



type Clock struct {
	elapsed_ns uint
	ns_cycle_length uint
}

func (c Clock) get_elapsed_ns() uint {
	return c.elapsed_ns
}

func (c Clock) reset() {
	c.elapsed_ns = 0
}

func (c Clock) tick() {
	c.elapsed_ns += c.ns_cycle_length;
}

func (c Clock) set_cycle_length(ns uint) {
	c.ns_cycle_length = ns
}