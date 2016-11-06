package iec62056

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/tarm/serial"
)

type (
	// Iec62056 represents a Iec62056-compatible meter.
	Iec62056 struct {
		port io.ReadWriteCloser
	}
)

var (
	// Start will mark the start of a message.
	Start = byte(0x2f)

	// End will mark the end of a message.
	End = byte(0x21)

	// FrameStart is used to mark the start of a frame.
	FrameStart = byte(0x02)

	// FrameEnd is used to mark the end of a frame.
	FrameEnd = byte(0x03)

	// LineFeed is a newline character in ascii.
	LineFeed = byte(0x0a)

	// Completion marks the end of multiple entities in the protocol.
	Completion = []byte{0x0d, LineFeed}

	// ErrAddressTooLong will be returned if the device address is too long.
	ErrAddressTooLong = errors.New("Device address too long (maximum is 32 characters)")
)

// NewIec62056 will initialize a new IEC-61107 reader with a user provided
// io.ReadWriteCloser.
func NewIec62056(port io.ReadWriteCloser) *Iec62056 {
	i := &Iec62056{
		port: port,
	}

	return i
}

// NewIec62056Serial will initilize a new reader for a IEC-62056-compatible
// meter.
func NewIec62056Serial(device string) (*Iec62056, error) {
	conf := &serial.Config{
		Name:        device,
		Baud:        300,
		Size:        7,
		Parity:      serial.ParityEven,
		ReadTimeout: time.Millisecond * 2000,
	}

	port, err := serial.OpenPort(conf)
	if err != nil {
		return nil, err
	}

	return NewIec62056(port), nil
}

// Close will close the connection to the meter.
func (i *Iec62056) Close() error {
	return i.port.Close()
}

func (i *Iec62056) read(length int, until *byte) ([]byte, error) {
	buf := make([]byte, 1024)
	var reply []byte

	for {
		n, err := i.port.Read(buf)
		if err != nil {
			return reply, err
		}

		if n > 0 {
			for _, b := range buf[0:n] {
				reply = append(reply, b)

				if until != nil && b == *until {
					return reply, nil
				}

				if len(reply) >= length {
					return reply, nil
				}
			}
		}
	}
}

// bcc will calculate a "block check character" according to ISO/IEC 1155:1978
func bcc(message []byte) byte {
	if len(message) < 1 {
		return 0
	}

	var reg byte

	// If the message doesn't start with a start marker, the start must be
	// included.
	if message[0] != FrameStart {
		reg ^= FrameStart
	}

	// xor everything baby :)
	for _, b := range message[1:] {
		reg ^= b
	}

	return reg
}

// Signin will start a session with the meter.
func (i *Iec62056) Signin(address string) ([]byte, ValueCollection, error) {
	if len(address) > 32 {
		return nil, nil, ErrAddressTooLong
	}

	// Say hello :)
	signin := fmt.Sprintf("/?%s!\r\n", address)
	i.port.Write([]byte(signin))

	// Read "identify" line
	identify, err := i.read(1000, &LineFeed)
	if err != nil {
		return nil, nil, err
	}

	// read message
	payload, err := i.read(3000, &FrameEnd)
	if err != nil {
		return identify, nil, err
	}

	// read and ignore checksum
	_, err = i.read(1, nil)
	if err != nil {
		return identify, nil, err
	}

	collection, err := NewValueCollection(payload)
	//	fmt.Printf("ID: \033[32m%v\033[0m\nPayload:\n\033[32m%v\033[0m\nBCC: 0x\033[32m%x\033[0m\nCalculated BCC: 0x\033[32m%x\033[0m\n", string(identify), string(payload[1:]), checksum[0], bcc(payload))
	return identify, collection, err
}
