package clocks

/*TestClock simulates the clock of a computer. Each instruction is executed, and the next one fetched,
according to the timing of the clock. This is a useful part of the simulation because it allows us to
speed up and slow down execution of the computer for observation and debugging.

The clock itself can be "turned off" so that execution proceeds as quickly as possible.
*/
type TestClock struct {
	elapsedNs     uint
	nsCycleLength uint
}

func (c *TestClock) getElapsedNs() uint {
	return c.elapsedNs
}

func (c *TestClock) reset() {
	c.elapsedNs = 0
}

func (c *TestClock) tick() {
	c.elapsedNs += c.nsCycleLength
}

func (c *TestClock) setCycleLength(ns uint) {
	c.nsCycleLength = ns
}
