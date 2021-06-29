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
	memory     memory
	escapeMode bool
	result     strings.Builder
}

func newFSM() *FSM {
	return &FSM{memory{}, false, strings.Builder{}}
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

func (s *FSM) done() error {
	if s.escapeMode {
		return ErrInvalidEscapedRune
	}
	if s.memory.used {
		s.result.WriteRune(s.memory.r)
	}
	return nil
}

func Unpack(input string) (string, error) {
	FSM := newFSM()
	for _, currentRune := range input {
		if FSM.escapeMode {
			if !escapeRune(currentRune) {
				return "", ErrInvalidEscapedRune
			}
			FSM.putInMemory(currentRune)
			FSM.escapeMode = false
		} else {
			switch {
			case currentRune == '\\':
				FSM.escapeMode = true
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
	if err := FSM.done(); err != nil {
		return "", err
	}
	return FSM.result.String(), nil
}
