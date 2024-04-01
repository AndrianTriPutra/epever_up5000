package epever

import (
	"context"
	"epever/atp/pkg/domain"
	"errors"
	"time"

	"github.com/simonvetter/modbus"
)

type strucT struct {
	setting Setting
}

type Setting struct {
	Port     string
	Baudrate uint
	Timeout  time.Duration
}

func NewRepository(setting Setting) RepositoryI {
	return &strucT{
		setting: setting,
	}
}

type RepositoryI interface {
	UP5000(ctx context.Context, id uint8) (data domain.Data, err error)
}

func (r strucT) UP5000(ctx context.Context, id uint8) (data domain.Data, err error) {
	url := "rtu://" + r.setting.Port
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      url,
		Speed:    r.setting.Baudrate,
		DataBits: 8,
		Parity:   modbus.PARITY_NONE,
		StopBits: 1,
		Timeout:  r.setting.Timeout,
	})
	if err != nil {
		err := errors.New("E0")
		return data, err
	}

	err = client.Open()
	if err != nil {
		err := errors.New("E1")
		return data, err
	}
	defer client.Close()

	client.SetUnitId(id)

	datas, err := client.ReadRegisters(0x3500, 5, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(0x3500, 5, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-1")
		return data, err
	}
	for i, value := range datas {
		switch i {
		case 0: //0x3500
			data.Input_Voltage = float32(value) / 100.00
		case 1: //0x3501
			data.Input_Current = float32(value) / 100.00

		case 4: //0x3504
			data.Input_Frequency = float32(value) / 100.00
		}
	}

	datas, err = client.ReadRegisters(0x3521, 2, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(0x3521, 2, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-2")
		return data, err
	}
	for i, value := range datas {
		switch i {
		case 0: //0x3521
			data.Output_Voltage = float32(value) / 100.00
		case 1: //0x3522
			data.Output_Current = float32(value) / 100.00

		}
	}

	dataS, err := client.ReadRegister(0x3414, modbus.INPUT_REGISTER)
	dataS, err = client.ReadRegister(0x3414, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-3")
		return data, err
	}
	data.Output_Frequency = float32(dataS) / 100.00

	datas, err = client.ReadRegisters(0x3580, 7, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(0x3580, 7, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-4")
		return data, err
	}
	for i, value := range datas {
		switch i {
		case 0: //0x3580
			data.Battery_Voltage = float32(value) / 100.00
		case 6: //0x3586
			data.Battery_SOC = float32(value)

		}
	}

	dataS, err = client.ReadRegister(0x3532, modbus.INPUT_REGISTER)
	dataS, err = client.ReadRegister(0x3532, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-5")
		return data, err
	}
	data.Inverter_Temperature = float32(dataS) / 100.00

	return data, nil
}
