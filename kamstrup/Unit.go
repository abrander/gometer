package kamstrup

import (
	"fmt"
)

type (
	// Unit will define an unit as known by Kamstrup meters.
	Unit byte
)

var (
	unitString = map[Unit]string{
		0x00: "",
		0x01: "Wh",
		0x02: "kWh",
		0x03: "MWh",
		0x04: "GWh",
		0x05: "j",
		0x06: "kj",
		0x07: "Mj",
		0x08: "Gj",
		0x09: "Cal",
		0x0a: "kCal",
		0x0b: "Mcal",
		0x0c: "Gcal",
		0x0d: "varh",
		0x0e: "kvarh",
		0x0f: "Mvarh",
		0x10: "Gvarh",
		0x11: "VAh",
		0x12: "kVAh",
		0x13: "MVAh",
		0x14: "GVAh",
		0x15: "kW",
		0x16: "kW",
		0x17: "MW",
		0x18: "GW",
		0x19: "kvar",
		0x1a: "kvar",
		0x1b: "Mvar",
		0x1c: "Gvar",
		0x1d: "VA",
		0x1e: "kVA",
		0x1f: "MVA",
		0x20: "GVA",
		0x21: "V",
		0x22: "A",
		0x23: "kV",
		0x24: "kA",
		0x25: "C",
		0x26: "K",
		0x27: "l",
		0x28: "m3",
		0x29: "l/h",
		0x2a: "m3/h",
		0x2b: "m3xC",
		0x2c: "ton",
		0x2d: "ton/h",
		0x2e: "h",
		0x2f: "hh:mm:ss",
		0x30: "yy:mm:dd",
		0x31: "yyyy:mm:dd",
		0x32: "mm:dd",
		0x33: " ",
		0x34: "bar",
		0x35: "RTC",
		0x36: "ASCII",
		0x37: "m3 x 10",
		0x38: "ton x 10",
		0x39: "GJ x 10",
		0x3a: "minutes",
		0x3b: "Bitfield",
		0x3c: "s",
		0x3d: "ms",
		0x3e: "days",
		0x3f: "RTC-Q",
		0x40: "Datetime",
	}
)

// String will convert the unit to a string.
func (u Unit) String() string {
	str := unitString[u]

	if str == "" {
		str = fmt.Sprintf("[UNKNOWN UNIT: %d [0x%x]", int(u), int(u))
	}

	return str
}
