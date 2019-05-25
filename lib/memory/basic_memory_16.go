package memory

import (
	"math"

	Utils "github.com/chenhowa/computer/lib/binaryInstructionExecution/bitUtils"
)

/*BasicMemory implements basic reading from and writing to bits in memory. It
supports a 16-bit address space and is byte-addressable. Reads are up to 32 bits,
as are writes
*/
type BasicMemory struct {
	memory     [65536]uint8
	maxAddress uint16
}

/*MakeBasicMemory constructs an instance of BasicMemory, where
all memory is default-initialized to 0*/
func MakeBasicMemory(maxAddress uint16) BasicMemory {
	memory := BasicMemory{
		maxAddress: maxAddress,
	}
	return memory
}

/*Get reads a doubleword from `address` in memory
If reading that doubleword involves reading from some invalid addresses,
the invalid addresses are regarded as having value `0`*/
func (m *BasicMemory) Get(address uint16) uint32 {
	if !(uint(address) < m.GetAddressSpaceSize()) {
		panic("Get: address was outside addressable memory")
	}

	var byteCount uint16
	var val uint32
	for ; byteCount < 4; byteCount++ {
		if uint(address)+uint(byteCount) < m.GetAddressSpaceSize() {
			val += uint32(m.memory[address+byteCount]) << (8 * byteCount)
		} else {
			break
		}
	}

	return val
}

/*Set writes the lowest `bitsToWrite` bits of `val` to memory at `address`. It starts writing
from the LSB of `val`, and it returns the number of bits successfully written. Attempting to
write more than 32 bits will write at most 32 bits. If memory runs out before `bitsToWrite` bits
are written, no explicit error occurs. The caller can detect this through the return value, which is the number
of bits successfully written*/
func (m *BasicMemory) Set(address uint16, val uint32, bitsToWrite uint) NumberOfBitsWritten {
	if !(uint(address) < m.GetAddressSpaceSize()) {
		panic("Set: address was outside addressable memory")
	}

	var bitsWritten NumberOfBitsWritten

	var trackedAddress = address
	for uint(bitsWritten) < bitsToWrite {
		if !(uint(trackedAddress) < m.GetAddressSpaceSize()) {
			break
		} else if bitsLeft := bitsToWrite - uint(bitsWritten); bitsLeft < 8 {
			keepTopBitsMask := uint8(Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), bitsLeft, 7))
			m.memory[trackedAddress] = (m.memory[trackedAddress] & keepTopBitsMask) |
				uint8(Utils.GetBitsInInclusiveRange(uint(val), uint(bitsWritten), uint(bitsWritten)+bitsLeft-1))
			bitsWritten += NumberOfBitsWritten(bitsLeft)
			break
		} else {
			m.memory[trackedAddress] = uint8(Utils.GetBitsInInclusiveRange(uint(val), uint(bitsWritten), uint(bitsWritten)+7))
			bitsWritten += 8
			trackedAddress++
		}
	}

	return bitsWritten
}

/*GetAddressSpaceSize returns the size of the address space. That is,
it returns <maximum valid address + 1>
*/
func (m *BasicMemory) GetAddressSpaceSize() uint {
	return uint(m.maxAddress) + 1
}
