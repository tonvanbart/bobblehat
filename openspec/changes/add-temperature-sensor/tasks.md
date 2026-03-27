## 1. Dependencies and Package Setup

- [ ] 1.1 Add `periph.io/x/conn/v3` to `go.mod`
- [ ] 1.2 Create `sense/sensors/` package directory with `sensors.go`

## 2. Core Implementation

- [ ] 2.1 Define LPS25HB constants (I2C address, register addresses, WHO_AM_I value)
- [ ] 2.2 Implement `Device` struct and `Open(bus i2c.Bus)` function with WHO_AM_I verification and sensor power-on
- [ ] 2.3 Implement `Temperature()` method — read registers `0x2B`/`0x2C`, convert with `int16(raw) / 480.0 + 42.5`
- [ ] 2.4 Implement `Close()` method

## 3. Testing

- [ ] 3.1 Write unit tests with a mock `i2c.Bus` verifying Open, Temperature, and Close behavior

## 4. Example

- [ ] 4.1 Create `examples/temperature/main.go` demonstrating temperature reading with `periph.io/x/host/v3`
