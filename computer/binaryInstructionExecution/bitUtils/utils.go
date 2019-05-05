package bitUtils

import (
	"fmt"
	"math"
)

/*signExtendUint32WithBit takes a uint32 `integer` and a `bit` number
(0 - 31 are the valid values) and copies the bit value at `bit` into
all the bits {n | 32 > n > `bit`} of `integer`, and then returns a copy
of that value
*/
func SignExtendUint32WithBit(integer uint32, bit uint) uint32 {
	bitValue := ((1 << bit) & integer) >> bit
	var mask uint32
	var signExtended uint32
	if max := math.MaxUint32; bitValue == 1 {
		mask = uint32(max << (bit + 1))
		signExtended = mask | integer
	} else {
		mask = uint32(max >> (32 - bit - 1))
		signExtended = mask & integer
	}

	return signExtended
}

/*KeepBitsInInclusiveRange takes `num` and zeros out all bits
not in the inclusive range from `start` to `end`

Example:
	KeepBitsInInclusiveRange(0b111, 1, 1) = 0b010
	KeepBitsInInclusiveRange(0b111, 1, 2) = 0b110
	KeepBitsInInclusiveRange(0b111, 0, 1) = 0b011
*/
func KeepBitsInInclusiveRange(num uint32, start uint, end uint) uint32 {
	if start > end {
		panic(fmt.Sprintf("keepBitsInRange: start > end, %d > %d", start, end))
	}

	majorMask := uint32(1<<(end+1)) - 1
	minorMask := uint32(math.MaxUint32) << (start)

	return (num & majorMask & minorMask)
}
