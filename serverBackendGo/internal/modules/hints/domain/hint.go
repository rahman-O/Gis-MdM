package domain

import (
	"errors"
	"strings"
)

// HintKey identifies a tutorial step shown to a user.
type HintKey string

const maxHintKeyLen = 100

var ErrEmptyHintKey = errors.New("error.hint.empty")

// ValidateHintKey checks non-empty and max length per DB column.
func ValidateHintKey(key string) (HintKey, error) {
	k := strings.TrimSpace(key)
	if k == "" {
		return "", ErrEmptyHintKey
	}
	if len(k) > maxHintKeyLen {
		return "", ErrEmptyHintKey
	}
	return HintKey(k), nil
}
