package kamstrup

import (
	"errors"
	"io"
	"time"

	"github.com/tarm/serial"
)

type (
	// Kamstrup represents a Kamstrup meter.
	Kamstrup struct {
		port *serial.Port
	}
)

var (
	// ErrWrongNumberOfRegisters will be returned if the meter by any chance
	// returns a wrong number of registers when replying to a GetRegister
	// command.
	ErrWrongNumberOfRegisters = errors.New("Wrong number of registers in reply")
)

// NewKamstrup will initilize a new Kamstrup. device should point to a serial
// device with an IR transceiver. Could be "/dev/ttyUSB0".
func NewKamstrup(device string) (*Kamstrup, error) {
	conf := &serial.Config{
		Name:        device,
		Baud:        9600,
		ReadTimeout: time.Millisecond * 200,
	}

	port, err := serial.OpenPort(conf)
	if err != nil {
		return nil, err
	}

	k := &Kamstrup{
		port: port,
	}

	return k, nil
}

// Close will close the connection to the meter.
func (k *Kamstrup) Close() error {
	return k.port.Close()
}

func (k *Kamstrup) readReply() ([]byte, error) {
	var out []byte

	// Read until "stop byte" or meter ack or EOF caused by a timeout.
readloop:
	for {
		buf := make([]byte, 128)

		n, err := k.port.Read(buf)
		if err == io.EOF {
			return out, nil
		}

		if err != nil {
			return nil, err
		}

		for i := 0; i < n; i++ {
			out = append(out, buf[i])

			if buf[i] == Stop {
				break readloop
			}

			if buf[i] == MeterAck {
				break readloop
			}
		}
	}

	return out, nil
}

// GetRegisters will read one or more values from the supplied registers.
func (k *Kamstrup) GetRegisters(registers ...uint16) (map[uint16]Value, error) {
	f := Frame{
		Type:      ToMeter,
		Address:   0x3f,
		CommandID: GetRegister,
	}

	// Build the GetRegister command payload.
	f.Data = make([]byte, 1, len(registers)+1)
	f.Data[0] = byte(len(registers))
	for _, register := range registers {
		f.Data = append(f.Data, byte(register>>8))
		f.Data = append(f.Data, byte(register&0xff))
	}

	reply, err := k.SendAndReceive(f)
	if err != nil {
		return nil, err
	}

	values := make(map[uint16]Value)
	pos := 0
	for r := 0; r < len(registers); r++ {
		if pos >= len(reply.Data) {
			break
		}

		reg := uint16(reply.Data[pos])<<8 + uint16(reply.Data[pos+1])
		pos += 2

		if pos >= len(reply.Data) {
			break
		}

		read, value, err := NewValue(reply.Data[pos:])
		if err == nil {
			values[reg] = value
		}

		pos += read
	}

	return values, nil
}

// GetRegister will read one value from the supplied register.
func (k *Kamstrup) GetRegister(register uint16) (Value, error) {
	results, err := k.GetRegisters(register)
	if err != nil {
		return Value{}, err
	}

	if len(results) != 1 {
		return Value{}, ErrWrongNumberOfRegisters
	}

	return results[register], nil
}

// SendAndReceive will send a frame and try to receive and decode a reply.
func (k *Kamstrup) SendAndReceive(frame Frame) (Frame, error) {
	var reply Frame

	payload := frame.Encode()
	_, err := k.port.Write(payload)
	if err != nil {
		return reply, err
	}

	raw, err := k.readReply()
	if err != nil {
		return reply, err
	}

	err = reply.Decode(raw)
	if err != nil {
		return reply, err
	}

	return reply, nil
}

// GetSerialNo will return the meter serial number.
func (k *Kamstrup) GetSerialNo() (int, error) {
	f := Frame{
		Type:      ToMeter,
		Address:   0x3f,
		CommandID: GetSerialNo,
	}

	reply, err := k.SendAndReceive(f)
	if err != nil {
		return 0, err
	}

	if len(reply.Data) < 4 {
		return 0, ErrFrameTooShort
	}

	sn := 0
	for i := 0; i < 4; i++ {
		sn <<= 8
		sn |= int(reply.Data[i])
	}

	return sn, nil
}

// GetType will return the type of meter. Please note that this is pretty
// arbitrary. Even the length varies between meters (!).
func (k *Kamstrup) GetType() ([]byte, error) {
	f := Frame{
		Type:      ToMeter,
		Address:   0x3f,
		CommandID: GetType,
	}

	reply, err := k.SendAndReceive(f)
	if err != nil {
		return []byte{}, err
	}

	return reply.Data, nil
}

// GetEventStatus will return the four event status bytes.
func (k *Kamstrup) GetEventStatus() (byte, byte, byte, byte, error) {
	f := Frame{
		Type:      ToMeter,
		Address:   0x3f,
		CommandID: GetType,
	}

	reply, err := k.SendAndReceive(f)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if len(reply.Data) != 4 {
		return 0, 0, 0, 0, ErrFrameTooShort
	}

	return reply.Data[0], reply.Data[1], reply.Data[2], reply.Data[3], nil
}
