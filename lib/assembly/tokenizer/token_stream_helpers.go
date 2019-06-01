package tokenizer

import "errors"

func (s *RiscVTokenStream) getCurrentChar() (byte, error) {
	if s.currentPosition < uint(len(s.input)) {
		return s.input[s.currentPosition], nil
	}

	return 0, errors.New("getCurrentChar: past end of input")
}

func (s *RiscVTokenStream) getNextChar() (byte, error) {
	nextPosition := s.currentPosition + 1
	if nextPosition < uint(len(s.input)) {
		return s.input[nextPosition], nil
	}

	return 0, errors.New("getNextChar: past end of input")
}

func (s *RiscVTokenStream) incrementCurrentPosition() {
	s.currentPosition++
}
