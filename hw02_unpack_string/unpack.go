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

type memory struct {
	r    rune
	used bool
}
type FSM struct {
	memory memory
	result strings.Builder
}

func newFSM() *FSM {
	return &FSM{memory{}, strings.Builder{}}
}

func (s *FSM) putInMemory(r rune) error {
	if s.memory.used {
		_, err := s.result.WriteRune(s.memory.r)
		if err != nil {
			return err
		}
	}
	s.memory = memory{r: r, used: true}
	return nil
}

func (s *FSM) repeat(r rune) error {
	if !s.memory.used {
		return ErrInvalidString
	}
	times := r - '0'
	for times > 0 {
		_, err := s.result.WriteRune(s.memory.r)
		if err != nil {
			return err
		}
		times--
	}
	s.memory.used = false
	return nil
}

func (s *FSM) done() string {
	if s.memory.used {
		s.result.WriteRune(s.memory.r)
	}
	return s.result.String()
}

func Unpack(input string) (string, error) {
	escapeMode := false
	FSM := newFSM()
	// add terminator zero rune
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
	return FSM.done(), nil
}
