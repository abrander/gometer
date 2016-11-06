package iec62056

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/abrander/gometer/kamstrup"
)

type (
	// ValueCollection is a collection of values.
	ValueCollection map[Obis]kamstrup.Value
)

func NewValueCollection(payload []byte) (ValueCollection, error) {
	c := make(ValueCollection)

	return c, c.ParsePayload(payload)
}

func parseValue(in string) (kamstrup.Value, error) {
	var value kamstrup.Value
	parts := strings.Split(in, "*")
	if len(parts) < 1 {
		return value, errors.New("Unable to parse value: " + in)
	}

	if len(parts) > 1 {
		value.Unit = kamstrup.UnitFromString(parts[1])
	}

	var err error
	value.Value, err = strconv.ParseFloat(parts[0], 64)
	//	v, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return value, errors.New("Unable to read value: " + in)
	}

	return value, nil
}

func parseLine(line string) (Obis, kamstrup.Value, error) {
	expr := regexp.MustCompile(`(?P<obis>[\d.]+)\((?P<value>.+)\)`)
	match := expr.FindStringSubmatch(line)

	result := make(map[string]string)
	for i, name := range expr.SubexpNames() {
		result[name] = match[i]
	}

	obis := NewObis(result["obis"])

	value, err := parseValue(result["value"])
	if err != nil {
		return obis, kamstrup.Value{}, err
	}

	return obis, value, nil
}

func (c ValueCollection) ParsePayload(payload []byte) error {
	// An empty payload is valid.
	if payload == nil {
		return nil
	}

	if bytes.HasPrefix(payload, []byte{FrameStart}) {
		payload = payload[1:]
	}

	if bytes.HasSuffix(payload, []byte{FrameEnd}) {
		payload = payload[:len(payload)-1]
	}

	lines := bytes.Split(payload, []byte{LineFeed})

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		str := string(line)
		str = strings.TrimSpace(str)

		obis, value, err := parseLine(str)
		if err == nil {
			c[obis] = value
		}
	}

	return nil
}
