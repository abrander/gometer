package kamstrup

import (
	"errors"
	"fmt"
	"math"
)

type (
	// Value is a value and a unit read from a Kamstrup meter.
	Value struct {
		Value float64
		Unit  Unit
	}
)

var (
	// ErrCouldNotDecodeValue will be returned from NewValue if the value
	// cannot be decoded.
	ErrCouldNotDecodeValue = errors.New("could not decode value")
)

// NewValue will initialize a new value based on raw bytes.
func NewValue(raw []byte) (int, Value, error) {
	value := Value{}

	l := len(raw)

	if l < 3 {
		return math.MaxInt64, value, ErrCouldNotDecodeValue
	}

	value.Unit = Unit(raw[0])
	mantissaLength := int(raw[1])

	if l < mantissaLength+3 {
		return math.MaxInt64, value, ErrCouldNotDecodeValue
	}

	mantissa := 0
	for i := 0; i < mantissaLength; i++ {
		mantissa <<= 8
		mantissa |= int(raw[i+3])
	}

	exponent := float64(raw[2] & 0x3f)
	if raw[2]&0x40 > 0 {
		exponent = -exponent
	}
	exponent = math.Pow(10, exponent)
	if raw[2]&0x80 > 0 {
		exponent = -exponent
	}

	value.Value = float64(mantissa) * exponent

	return mantissaLength + 3, value, nil
}

// String will return a string representation of the value and unit.
func (v Value) String() string {
	return fmt.Sprintf("%.3f %s", v.Value, v.Unit.String())
}
