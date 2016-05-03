package kamstrup

import "errors"

type (
	// Frame represents a frame sent to/or received from a Kamstrup meter (or
	// similar device).
	Frame struct {
		Type      byte
		Address   byte
		CommandID byte
		Data      []byte
	}
)

const (
	// ToMeter is used as start byte in frames to the meter.
	ToMeter = byte(0x80)

	// FromMeter is used as start byte by the meter.
	FromMeter = byte(0x40)

	// MeterAck is used by the meter to acknowledge a command.
	MeterAck = byte(0x06)

	// Stop indicates the end of a frame.
	Stop = byte(0x0d)
)

var (
	// Byte values which must be escaped before transmission.
	escapes = map[byte]bool{
		Stop:      true,
		MeterAck:  true,
		0x1b:      true,
		FromMeter: true,
		ToMeter:   true,
	}

	// ErrInvalidChecksum means that we were unable to verify the checksum from
	// the meter.
	ErrInvalidChecksum = errors.New("checksum did not validate")

	// ErrFrameEmpty will be returned if the frame is empty.
	ErrFrameEmpty = errors.New("frame is empty")

	// ErrFrameTooShort will be returned if we tries to decode a frame that
	// looks too short.
	ErrFrameTooShort = errors.New("frame too short")

	// ErrInvalidFrame will be returned if we do not recognize the frame type.
	ErrInvalidFrame = errors.New("invalid frame type")
)

// Decode will try to decode a raw frame from the wire.
func (f *Frame) Decode(raw []byte) error {
	frameLength := len(raw)

	if frameLength == 0 {
		return ErrFrameEmpty
	}

	var unescaped []byte
	unescaped = append(unescaped, raw[0])
	for i := 1; i < frameLength-1; i++ {
		b := raw[i]

		if b == 0x1b {
			b = raw[i+1] ^ 0xff
			i++
		}

		unescaped = append(unescaped, b)
	}
	unescaped = append(unescaped, raw[frameLength-1])
	frameLength = len(unescaped)

	switch unescaped[0] {
	case MeterAck:
		f.Type = MeterAck
		return nil
	case 0x0:
		fallthrough
	case ToMeter:
		fallthrough
	case FromMeter:
		if frameLength < 6 {
			return ErrFrameTooShort
		}
	default:
		return ErrInvalidFrame
	}

	f.Type = unescaped[0]
	f.Address = unescaped[1]
	f.CommandID = unescaped[2]
	f.Data = unescaped[3 : frameLength-3]

	checksum := f.checksum()

	// Check if we got a stop byte.
	if unescaped[frameLength-1] != Stop {
		return ErrFrameTooShort
	}

	// Verify the checksum.
	if byte(checksum>>8) != unescaped[frameLength-3] || byte(checksum&0xff) != unescaped[frameLength-2] {
		return ErrInvalidChecksum
	}

	return nil
}

// Encode will encode a complete frame including start and stop bytes ready for the wire.
func (f Frame) Encode() []byte {
	var payload []byte

	payload = append(payload, f.Address)
	payload = append(payload, f.CommandID)
	payload = append(payload, f.Data...)

	// Append these two fields to compute checksum.
	payload = append(payload, 0x0)
	payload = append(payload, 0x0)

	// Make checksum of everything but direction byte.
	checksum := f.checksum()

	// Replace the two "blank" checksums bytes with the real checksum.
	payload[len(payload)-2] = byte(checksum >> 8)
	payload[len(payload)-1] = byte(checksum & 0xff)

	// Escape and construct wire bytes.
	var raw []byte
	raw = append(raw, f.Type)
	for _, b := range payload {
		// Check if we need escaping
		if escapes[b] {
			raw = append(raw, 0x1b)
			raw = append(raw, b^0xff)
		} else {
			raw = append(raw, b)
		}
	}
	raw = append(raw, Stop)

	return raw
}

func (f Frame) checksum() int {
	var msg []byte

	msg = append(msg, f.Address)
	msg = append(msg, f.CommandID)
	msg = append(msg, f.Data...)
	msg = append(msg, 0x0)
	msg = append(msg, 0x0)

	// Kamstrup almost uses CRC16-CCITT ...
	poly := 0x1021

	// ... but with an initial value of 0.
	reg := 0x0000

	for _, b := range msg {
		mask := byte(0x80)

		for mask > 0 {
			reg <<= 1

			if b&mask > 0 {
				reg |= 1
			}

			mask >>= 1

			if reg&0x10000 > 0 {
				reg &= 0xffff
				reg ^= poly
			}
		}
	}

	return reg
}
