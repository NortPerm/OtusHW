package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidString      = errors.New("invalid string")
	ErrInvalidEscapedRune = errors.New("invalid ecsape-rune")
)

func escapeRune(r rune) bool {
	return r == '\\' || unicode.IsDigit(r)
}

type FSM struct {
	memory rune
	result strings.Builder
}

func newFSM() *FSM {
	return &FSM{0, strings.Builder{}}
}

func (s *FSM) putInMemory(r rune) error {
	if s.memory > 0 {
		_, err := s.result.WriteRune(s.memory)
		if err != nil {
			return err
		}
	}
	s.memory = r
	return nil
}

func (s *FSM) repeat(r rune) error {
	if s.memory == 0 {
		return ErrInvalidString
	}
	times := r - '0'
	for times > 0 {
		_, err := s.result.WriteRune(s.memory)
		if err != nil {
			return err
		}
		times--
	}
	s.memory = 0
	return nil
}

func Unpack(input string) (string, error) {
	escapeMode := false
	FSM := newFSM()
	// add terminator zero rune
	input += "\x00"
	for _, currentRune := range input {
		if escapeMode {
			if !escapeRune(currentRune) {
				return "", ErrInvalidEscapedRune
			}
			FSM.putInMemory(currentRune)
			escapeMode = false
		} else {
			switch {
			case currentRune == '\\':
				escapeMode = true
			case unicode.IsDigit(currentRune):
				if err := FSM.repeat(currentRune); err != nil {
					return "", err
				}
			default:
				if err := FSM.putInMemory(currentRune); err != nil {
					return "", err
				}
			}
		}
	}
	return FSM.result.String(), nil
}
