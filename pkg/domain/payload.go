package domain

type Payload struct {
	Device_ID string
	Timestamp string
	Inverter  Inverter
}

type Inverter struct {
	ModbusID uint8
	Data     Data
}

type Data struct {
	Battery_Voltage float32
	Battery_SOC     float32

	Input_Voltage   float32
	Input_Current   float32
	Input_Frequency float32

	Output_Voltage   float32
	Output_Current   float32
	Output_Frequency float32

	Inverter_Temperature float32
}
