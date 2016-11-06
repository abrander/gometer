package iec62056

import (
	"strings"
)

type (
	// Obis is the iec62056 "OBject Identification System"
	Obis struct {
		A, B, C, D, E, F string
	}
)

var (
	aDescription = map[string]string{
		"0": "Abstract objects",
		"1": "Electricity",
		"4": "Heating costs",
		"5": "Cooling energy",
		"6": "Heat",
		"7": "Gas",
		"8": "Cold Water",
		"9": "Warm water",
	}

	bDescription = map[string]string{}

	cDescription = map[string]string{
		"0":  "General purpose objects",
		"1":  "Σ Li active power + (import)",
		"2":  "Σ Li active power - (export )",
		"3":  "Σ Li reactive power +",
		"4":  "Σ Li reactive power -",
		"5":  "Σ Li reactive power Q I",
		"6":  "Σ Li reactive power Q II",
		"7":  "Σ Li reactive power Q III",
		"8":  "Σ Li reactive power Q IV",
		"9":  "Σ Li apparent power +",
		"10": "Σ Li apparent power -",
		"11": "",
		"12": "",
		"13": "Power factor",
		"14": "Frequency",

		"21": "L1 active power +",
		"22": "L1 active power -",
		"23": "L1 reactive power +",
		"24": "L1 reactive power -",
		"25": "L1 reactive power Q I",
		"26": "L1 reactive power Q II",
		"27": "L1 reactive power Q III",
		"28": "L1 reactive power Q IV",
		"29": "L1 apparent power +",
		"30": "L1 apparent power -",
		"31": "L1 phase A current",
		"32": "L1 phase A voltage",
		"33": "L1 power factor",
		"41": "L2 active power +",
		"42": "L2 active power -",
		"43": "L2 reactive power +",
		"44": "L2 reactive power -",
		"45": "L2 reactive power Q I",
		"46": "L2 reactive power Q II",
		"47": "L2 reactive power Q III",
		"48": "L2 reactive power Q IV",
		"49": "L2 apparent power +",
		"50": "L2 apparent power -",
		"51": "L2 phase B current",
		"52": "L2 phase B voltage",
		"53": "L2 power factor",
		"61": "L3 active power +",
		"62": "L3 active power -",
		"63": "L3 reactive power +",
		"64": "L3 reactive power -",
		"65": "L3 reactive power Q I",
		"66": "L3 reactive power Q II",
		"67": "L3 reactive power Q III",
		"68": "L3 reactive power Q IV",
		"69": "L3 apparent power +",
		"70": "L3 apparent power -",
		"71": "L3 phase C current",
		"72": "L3 phase C voltage",
		"73": "L3 power factor",
		"94": "Country specific OBIS codes for different countries",
		"C":  "Service information",
		"F":  "Error messages",
		"L":  "List objects",
		"P":  "Data profiles",
	}

	dDescription = map[string]string{
		"0":  "",
		"1":  "Cumulative minimum 1",
		"2":  "Cumulative maximum 1",
		"3":  "Minimum 1",
		"4":  "Average value 1 current measuring period",
		"5":  "Average value 1 last measuring period",
		"6":  "Maximum 1",
		"7":  "Instantaneous value",
		"8":  "Energy",
		"9":  "Time integral 2",
		"10": "Time integral 3",
		"11": "Cumulative minimum 2",
		"12": "Cumulative maximum 2",
		"13": "Minimum 2",
		"14": "Average value 2 current measuring period",
		"15": "Average value 2 last measuring period",
		"16": "Maximum 2",
		"21": "Cumulative minimum 3",
		"22": "Cumulative maximum 3",
		"23": "Minimum 3",
		"24": "Average value 3 current measuring period",
		"25": "Average value 3 last measuring period",
		"26": "Maximum 3",
		"27": "",
		"28": "",
		"29": "Energy feed",
		"55": "Test equipment",
		"58": "Testing time integral",
		"F":  "Error message ",
	}

	eDescription = map[string]string{
		"0": "Total",
		"1": "Tariff 1",
		"2": "Tariff 2",
		"3": "Tariff 3",
		"4": "Tariff 4",
		"5": "Tariff 5",
		"6": "Tariff 6",
		"7": "Tariff 7",
		"8": "Tariff 8",
		"9": "Tariff 9",
	}

	fDescription = map[string]string{}
)

// NewObis will istantiate a new Obis.
func NewObis(raw string) Obis {
	o := Obis{}

	o.Parse(raw)

	return o
}

// Parse will try to parse a raw string as an OBIS code. Please note that this
// will do no error checking at all.
func (o *Obis) Parse(raw string) {
	var current *string

	// If we see a ":" in the raw string, we start at A.
	if strings.Contains(raw, ":") {
		current = &o.A
	} else {
		current = &o.C
	}

	for _, b := range []byte(raw) {
		switch b {
		case byte(0x2d): // -
			// A ended.
			current = &o.B
		case byte(0x3a): // :
			// A or B ended.
			current = &o.C
		case byte(0x2e): // .
			// C or D ended.
			switch current {
			case &o.C:
				current = &o.D
			case &o.D:
				current = &o.E
			}
		case byte(0x2a): // *
			fallthrough
		case byte(0x26): // &
			// We're at F.
			current = &o.F
		default:
			// It's not a "control character", it must be part of a value.
			*current += string(b)
		}
	}
}

// Description will return a human readable description of the OBIS code (if
// available).
func (o *Obis) Description() string {
	description := ""

	if aDescription[o.A] != "" {
		description += aDescription[o.A]
	}

	if bDescription[o.B] != "" {
		description += " " + bDescription[o.B]
	}

	if cDescription[o.C] != "" {
		description += " " + cDescription[o.C]
	}

	if dDescription[o.D] != "" {
		description += " " + dDescription[o.D]
	}

	if eDescription[o.E] != "" {
		description += " " + eDescription[o.E]
	}

	if fDescription[o.F] != "" {
		description += " " + fDescription[o.F]
	}

	return strings.TrimSpace(description)
}
