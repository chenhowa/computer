package instructionParsing

/*BinaryInstructionBuilder builds instructions starting
from the LSB and going to the MSB. It publishes the resultant
instruction as a uint, which will have to be cast to the correct
size.
*/
type BinaryInstructionBuilder struct {
	maxBits      uint
	currentBit   uint
	currentValue uint
}

/*MakeInstructionBuilder constructs a binaryInstructionBuilder
with `maxBits` bits. General rule of thumb is that the
number of bits should be 64 or lower.
*/
func MakeInstructionBuilder(maxBits uint) BinaryInstructionBuilder {
	builder := BinaryInstructionBuilder{
		maxBits:      maxBits,
		currentBit:   0,
		currentValue: 0,
	}
	return builder
}

/*AddNextXBits adds the `x` lowest bits to the next `x` bits of
the instruction that will be generated. For example, on the first call,
calling builder.AddNextXBits(2, 5)

However, if adding the next `x` bits would exceed the `builder.maxBits`
count, then the most significant excess bits will be truncated.
*/
func (builder *BinaryInstructionBuilder) AddNextXBits(x uint, value uint) {
	var numBits = x
	if (builder.currentBit + x) > builder.maxBits {
		numBits = builder.maxBits - builder.currentBit
	}
	mask := uint((1 << numBits) - 1)
	lowestXBits := mask & value
	shiftedBits := lowestXBits << builder.currentBit
	builder.currentValue += shiftedBits
	builder.currentBit += numBits
}

/*Build function builds the instruction and
publishes it as a uint.
*/
func (builder *BinaryInstructionBuilder) Build() uint {
	return builder.currentValue
}
