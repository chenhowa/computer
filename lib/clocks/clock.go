package clocks

/*Clock is a struct that simulates the clock in an actual CPU.
The Clock can be tuned to "tick" based on a time interval, or based on external input
*/
type Clock struct {
	waiter waiter
	count  uint
}

/*MakeClock constructs an instance of Clock*/
func MakeClock(waiter waiter) Clock {
	clock := Clock{
		waiter: waiter,
		count:  0,
	}

	return clock
}

type waiter interface {
	wait()
}

func (c *Clock) waitForNextCycle() {
	c.waiter.wait()
	c.count += 1
}

func (c *Clock) reset() {
	c.count = 0
}
