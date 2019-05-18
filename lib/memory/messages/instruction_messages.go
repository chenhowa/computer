package messages

/*InstructionMessages has methods for returning messages to give to command line user
when requesting instructions from them*/
type InstructionMessages struct{}

func (m *InstructionMessages) getInputPromptMessage(address string, value string) string {
	return "Instruction " + value + " at address: " + address + "\n" +
		"Press <Enter> to keep this instruction, or enter your own: "
}

func (m *InstructionMessages) getInputErrorMessage() string {
	return "Invalid instruction. Please enter a valid instruction according to the RiscV assembly language"
}
