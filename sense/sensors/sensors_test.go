package sensors

import (
	"encoding/binary"
	"errors"
	"math"
	"testing"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
)

// mockBus implements i2c.Bus for testing.
type mockBus struct {
	txFunc func(addr uint16, w, r []byte) error
}

func (m *mockBus) Tx(addr uint16, w, r []byte) error {
	return m.txFunc(addr, w, r)
}

func (m *mockBus) SetSpeed(freq physic.Frequency) error { return nil }

func (m *mockBus) String() string { return "mock" }

// lps25hbMock creates a mock bus that simulates a working LPS25HB sensor.
// tempRaw is the 16-bit raw temperature value the sensor returns.
func lps25hbMock(tempRaw int16) *mockBus {
	return &mockBus{
		txFunc: func(addr uint16, w, r []byte) error {
			if addr != lps25hbAddr {
				return errors.New("unexpected address")
			}
			if len(w) == 0 {
				return errors.New("no register specified")
			}
			switch w[0] {
			case regWhoAmI:
				if len(r) > 0 {
					r[0] = lps25hbWhoAmI
				}
			case regCtrl1:
				// power-on write, nothing to return
			case regTempOutL | 0x80:
				if len(r) >= 2 {
					binary.LittleEndian.PutUint16(r, uint16(tempRaw))
				}
			default:
				return errors.New("unexpected register")
			}
			return nil
		},
	}
}

func TestOpen(t *testing.T) {
	bus := lps25hbMock(0)
	dev, err := Open(bus)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	if dev == nil {
		t.Fatal("Open() returned nil device")
	}
}

func TestOpenWrongWhoAmI(t *testing.T) {
	bus := &mockBus{
		txFunc: func(addr uint16, w, r []byte) error {
			if w[0] == regWhoAmI && len(r) > 0 {
				r[0] = 0x00 // wrong ID
			}
			return nil
		},
	}
	_, err := Open(bus)
	if err == nil {
		t.Fatal("Open() should fail with wrong WHO_AM_I")
	}
}

func TestOpenI2CError(t *testing.T) {
	bus := &mockBus{
		txFunc: func(addr uint16, w, r []byte) error {
			return errors.New("i2c error")
		},
	}
	_, err := Open(bus)
	if err == nil {
		t.Fatal("Open() should fail on I2C error")
	}
}

func TestTemperature(t *testing.T) {
	tests := []struct {
		name  string
		raw   int16
		wantC float64
	}{
		{"zero raw", 0, 42.5},
		{"room temp ~22°C", -9840, 22.0},
		{"positive raw", 480, 43.5},
		{"negative temp", -20400, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := lps25hbMock(tt.raw)
			dev, err := Open(bus)
			if err != nil {
				t.Fatalf("Open() error: %v", err)
			}
			got, err := dev.Temperature()
			if err != nil {
				t.Fatalf("Temperature() error: %v", err)
			}
			if math.Abs(got-tt.wantC) > 0.01 {
				t.Errorf("Temperature() = %f, want %f", got, tt.wantC)
			}
		})
	}
}

func TestTemperatureAfterClose(t *testing.T) {
	bus := lps25hbMock(0)
	dev, err := Open(bus)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	dev.Close()
	_, err = dev.Temperature()
	if err == nil {
		t.Fatal("Temperature() should fail after Close()")
	}
}

// Verify mockBus implements i2c.Bus.
var _ i2c.Bus = (*mockBus)(nil)
