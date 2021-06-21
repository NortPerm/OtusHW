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

type scanner struct {
	memory rune
	result strings.Builder
}

func newScanner() *scanner {
	return &scanner{0, strings.Builder{}}
}

func (s *scanner) putInMemory(r rune) error {
	if s.memory > 0 {
		_, err := s.result.WriteRune(s.memory)
		if err != nil {
			return err
		}
	}
	s.memory = r
	return nil
}

func (s *scanner) repeat(r rune) error {
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
	scan := newScanner()
	// add terminator zero rune
	input += "\x00"
	for _, currentRune := range input {
		if escapeMode {
			if !escapeRune(currentRune) {
				return "", ErrInvalidEscapedRune
			}
			scan.putInMemory(currentRune)
			escapeMode = false
		} else {
			switch {
			case currentRune == '\\':
				escapeMode = true
			case unicode.IsDigit(currentRune):
				if err := scan.repeat(currentRune); err != nil {
					return "", err
				}
			default:
				if err := scan.putInMemory(currentRune); err != nil {
					return "", err
				}
			}
		}
	}
	return scan.result.String(), nil
}
