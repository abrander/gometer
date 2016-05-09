package kamstrup

// Known addresses for Kamstrup 382 kWh meters.
const (
	EnergyIn       = uint16(0x0001) // "Energy in"
	EnergyOut      = uint16(0x0002) // "Energy out"
	EnergyInHiRes  = uint16(0x000d) // "Energy in hi-res"
	EnergyOutHiRes = uint16(0x000e) // "Energy out hi-res"

	VoltageP1 = uint16(0x041e) // "Voltage p1"
	VoltageP2 = uint16(0x041f) // "Voltage p2"
	VoltageP3 = uint16(0x0420) // "Voltage p3"

	CurrentP1 = uint16(0x0434) // "Current p1"
	CurrentP2 = uint16(0x0435) // "Current p2"
	CurrentP3 = uint16(0x0436) // "Current p3"

	InternalTemperature = uint16(0x0437) // Internal meter temperature.

	PowerP1 = uint16(0x0438) // "Power p1"
	PowerP2 = uint16(0x0439) // "Power p2"
	PowerP3 = uint16(0x043a) // "Power p3"
)

// Known registers for Multical 601.
const (
	Date    = uint16(0x03eb) // Current date (YYMMDD)
	Energy1 = uint16(0x003c) // Energy register 1: Heat energy
	Energy2 = uint16(0x005e) // Energy register 2: Control energy
	Energy3 = uint16(0x003f) // Energy register 3: Cooling energy
	Energy4 = uint16(0x003d) // Energy register 4: Flow energy
	Energy5 = uint16(0x003e) // Energy register 5: Return flow energy
	Energy6 = uint16(0x005f) // Energy register 6: Tap water energy
	Energy7 = uint16(0x0060) // Energy register 7: Heat energy Y
	Energy8 = uint16(0x0061) // Energy register 8: [m³ * T1]
	Energy9 = uint16(0x006e) // Energy register 9: [m³ * T2]
)

// Useful aliases for Multical 601 registers.
const (
	HeatEnergy       = Energy1
	ControlEnergy    = Energy2
	CoolingEnergy    = Energy3
	FlowEnergy       = Energy4
	ReturnFlowEnergy = Energy5
	TapWaterEnergy   = Energy6
	HeatEnergyY      = Energy7
)
