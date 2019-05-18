package delay

import (
	"time"
)

/*TimeDelay will delay for a given amount of time by waiting on a go channel.
 */
type TimeDelay struct {
}

/*DelayForMilliseconds pauses the executing goroutine for specified number of `milliseconds`
 */
func (d *TimeDelay) DelayForMilliseconds(milliseconds uint) {
	nanoseconds := milliseconds * 1000 * 1000
	time.Sleep(time.Duration(nanoseconds))
}
