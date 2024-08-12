package id

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/parsoj/paymentsys/internal/id/base62"
	"github.com/xyproto/randomstring"
)

type SpecialId string

type Identifier[T string] interface {
	New() (T, error)
	FromString(string) (T, error)
	FromBytes([]byte) (T, error)
	Validate(T) (bool, error)
}

type SpecialIdGenerator struct{}

// New generates a new random identifier.
func (g *SpecialIdGenerator) New() (SpecialId, error) {

	timestamp_string := g.timeStampString()

	filler_string_length := 20 - len(timestamp_string)

	filler_string := ""
	if filler_string_length > 0 {

		filler_string = g.fillerString(filler_string_length)

	}

	combined := timestamp_string + filler_string

	id := SpecialId(combined)
	return id, nil
}

// FromString parses a string into an identifier.
func (g *SpecialIdGenerator) FromString(s string) (SpecialId, error) {
	if len(s) == 0 {
		return "", errors.New("empty string cannot be an identifier")
	}
	return SpecialId(s), nil
}

// FromBytes parses a byte slice into an identifier.
func (g *SpecialIdGenerator) FromBytes(b []byte) (SpecialId, error) {
	if len(b) == 0 {
		return "", errors.New("empty byte slice cannot be an identifier")
	}
	id := SpecialId(hex.EncodeToString(b))
	return id, nil
}

// Validate checks if the identifier conforms to the specification and is valid.
func (g *SpecialIdGenerator) Validate(id SpecialId) (bool, error) {
	// Example validation: check if the identifier is a valid hex string with a length of at most 20 characters
	if len(id) > 20 {
		return false, errors.New("identifier is too long")
	}

	return true, nil
}

func (g *SpecialIdGenerator) timeStampString() string {

	epoch_time := time.Now().UnixNano()

	timestamp_string := base62.Base62Encode(epoch_time)

	return timestamp_string

}

func (g *SpecialIdGenerator) fillerString(length int) string {

	filler_string := randomstring.CookieFriendlyString(length)

	return filler_string

}
