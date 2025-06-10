package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const shield rune = 92

var ErrInvalidString = errors.New("invalid string")

func checkDigit(curIdx int, curValue rune, prevValue rune, prevIsShielded bool) (int, bool, error) {
	count := 0
	if curIdx == 0 {
		return count, prevIsShielded, ErrInvalidString
	}
	if unicode.IsDigit(prevValue) && !prevIsShielded {
		return count, prevIsShielded, ErrInvalidString
	}
	if prevValue == shield && !prevIsShielded {
		prevIsShielded = true
	} else {
		var err error
		count, err = strconv.Atoi(string(curValue))
		if err != nil {
			return count, prevIsShielded, ErrInvalidString
		}
		prevIsShielded = false
	}
	return count, prevIsShielded, nil
}

func checkRune(curValue rune, prevValue rune, prevIsShielded bool) (int, bool, error) {
	count := 0
	if prevValue == shield && !prevIsShielded && curValue != shield {
		return count, prevIsShielded, ErrInvalidString
	}
	if prevValue == 0 {
		return count, prevIsShielded, nil
	}
	if prevIsShielded {
		count = 1
		prevIsShielded = false
	} else {
		if prevValue != shield && !unicode.IsDigit(prevValue) {
			count = 1
		} else {
			prevIsShielded = true
		}
	}
	return count, prevIsShielded, nil
}

func Unpack(s string) (string, error) {
	result := &strings.Builder{}
	var prevValue rune
	prevIsShielded := false

	for key, value := range s {
		repeatCount := 0
		if unicode.IsDigit(value) {
			var err error
			repeatCount, prevIsShielded, err = checkDigit(key, value, prevValue, prevIsShielded)
			if err != nil {
				return "", ErrInvalidString
			}
		} else {
			var err error
			repeatCount, prevIsShielded, err = checkRune(value, prevValue, prevIsShielded)
			if err != nil {
				return "", ErrInvalidString
			}
		}
		result.WriteString(strings.Repeat(string(prevValue), repeatCount))
		prevValue = value
	}
	if prevValue != 0 && !unicode.IsDigit(prevValue) ||
		(unicode.IsDigit(prevValue) && prevIsShielded) {
		result.WriteString(string(prevValue))
	}
	return result.String(), nil
}
