package kamstrup

// Known commands available on Kamstrup meters.
const (
	GetType           = byte(0x01) // This command returns identification of the type of meter and software revision.
	GetSerialNo       = byte(0x02) // This command returns the serial number of the meter.
	SetClock          = byte(0x09) // This command sets the meter clock.
	GetRegister       = byte(0x10) // This command returns a variable set of registers.
	PutRegister       = byte(0x11) // This command is used for changing the value of a given register.
	GetEventStatus    = byte(0x9b) // Returns the four event status bytes.
	ClearEventStatus  = byte(0x9c) // Clear the event status bytes.
	GetLogTimePresent = byte(0xa0) // Log readout from specified ‘timestamp’ towards ‘now’.
	GetLogLastPresen  = byte(0xa1) // Log readout from ‘LRORecID’ towards ‘now’.
	GetLogIDPresent   = byte(0xa2) // Log readout from specified ‘Record ID Requested’ towards ‘now’.
	GetLogTimePast    = byte(0xa3) // Log readout from specified ‘timestamp’ towards the ‘past’.
)
