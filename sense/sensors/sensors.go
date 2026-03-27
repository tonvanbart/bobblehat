// Package sensors provides access to the Sense HAT's onboard sensors.
package sensors

import (
	"encoding/binary"
	"errors"
	"fmt"

	"periph.io/x/conn/v3/i2c"
)

// LPS25HB constants
const (
	lps25hbAddr    uint16 = 0x5C
	lps25hbWhoAmI  byte   = 0xBD
	regWhoAmI      byte   = 0x0F
	regCtrl1       byte   = 0x20
	regTempOutL    byte   = 0x2B
	regTempOutH    byte   = 0x2C
	ctrlPowerOnODR byte   = 0x84 // PD=1 (power on), ODR=100 (1 Hz)
)

// Device represents a connection to the Sense HAT sensors.
type Device struct {
	dev    *i2c.Dev
	closed bool
}

// Open initializes the LPS25HB sensor on the given I2C bus and returns a Device.
// It verifies the sensor is present via WHO_AM_I and powers it on.
func Open(bus i2c.Bus) (*Device, error) {
	dev := &i2c.Dev{Bus: bus, Addr: lps25hbAddr}

	// Verify WHO_AM_I
	var id [1]byte
	if err := dev.Tx([]byte{regWhoAmI}, id[:]); err != nil {
		return nil, fmt.Errorf("sensors: failed to read WHO_AM_I: %w", err)
	}
	if id[0] != lps25hbWhoAmI {
		return nil, fmt.Errorf("sensors: unexpected WHO_AM_I: got 0x%02X, want 0x%02X", id[0], lps25hbWhoAmI)
	}

	// Power on with 1 Hz output data rate
	if err := dev.Tx([]byte{regCtrl1, ctrlPowerOnODR}, nil); err != nil {
		return nil, fmt.Errorf("sensors: failed to power on LPS25HB: %w", err)
	}

	return &Device{dev: dev}, nil
}

// Temperature returns the current temperature in degrees Celsius.
func (d *Device) Temperature() (float64, error) {
	if d.closed {
		return 0, errClosed
	}

	var buf [2]byte
	// Read TEMP_OUT_L and TEMP_OUT_H with auto-increment (MSB of register address set)
	if err := d.dev.Tx([]byte{regTempOutL | 0x80}, buf[:]); err != nil {
		return 0, fmt.Errorf("sensors: failed to read temperature: %w", err)
	}

	raw := int16(binary.LittleEndian.Uint16(buf[:]))
	return float64(raw)/480.0 + 42.5, nil
}

// Close releases the device. Subsequent calls to Temperature return an error.
func (d *Device) Close() error {
	d.closed = true
	return nil
}

var errClosed = errors.New("sensors: device is closed")
