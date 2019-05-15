package sources

import (
	"fmt"
)

/*CommandLineSource is a source that takes the value from the user at
the command line*/
type CommandLineSource struct {
	inputTransformer  inputTransformer
	outputTransformer outputTransformer
	messageSource     messageSource
}

/*MakeCommandLineSource is a constructor for CommandLineSource*/
func MakeCommandLineSource(transformer inputTransformer) CommandLineSource {
	source := CommandLineSource{
		inputTransformer: transformer,
	}

	return source
}

type inputTransformer interface {
	transform(input string) (uint32, error)
}

type outputTransformer interface {
	transform(output uint32) string
}

type messageSource interface {
	getInputErrorMessage() string
	getInputPromptMessage(address uint16, value string) string
}

/*Get prints a prompt to stdout asking whether the user wants to accept
the existing value for this address, or to enter their own. If they want to enter their own,
they can simply press <Enter>. If they want to enter their own, they can write their own,
and press <Enter> to submit it. If the value they entered is valid, then it is returned, otherwise
the user will be prompted again until they enter a valid value.
*/
func (s *CommandLineSource) Get(address uint16, existingVal uint32) uint32 {
	var input string
	var result uint32
	done := false

	inputMessage := s.messageSource.getInputPromptMessage(address, s.outputTransformer.transform(existingVal))
	errorMessage := s.messageSource.getInputErrorMessage()

	for !done {
		fmt.Printf(inputMessage)
		_, readErr := fmt.Scanln(&input)
		transformed, tErr := s.inputTransformer.transform(input)
		if (readErr != nil) || (tErr != nil) {
			fmt.Printf(errorMessage)
		} else {
			result = transformed
			done = true
		}
	}

	return result
}
