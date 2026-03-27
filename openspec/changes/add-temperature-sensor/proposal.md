## Why

The Sense HAT has onboard weather sensors (pressure, humidity, temperature) but bobblehat currently only supports the LED matrix and joystick. Adding temperature sensor support is the first step toward full weather sensor coverage, and provides a foundation (I2C access pattern, package structure) for pressure and humidity to follow.

## What Changes

- Add a new `sense/sensors` package that provides a high-level API for reading Sense HAT sensor data
- Implement temperature reading from the LPS25HB pressure/temperature sensor (I2C address `0x5C`)
- Introduce `periph.io/x/conn/v3/i2c` as a dependency for I2C bus interfaces
- Add an example program demonstrating temperature readings

## Capabilities

### New Capabilities
- `temperature-reading`: Read ambient temperature in °C from the LPS25HB sensor via I2C, exposed through a `sense/sensors` package API

### Modified Capabilities

(none)

## Impact

- **New package**: `sense/sensors` — new public API surface
- **New dependency**: `periph.io/x/conn/v3` (interface-only; applications must import `periph.io/x/host/v3` to get real I2C drivers)
- **Hardware**: Requires I2C bus access (`/dev/i2c-1`) on the Raspberry Pi
- **No breaking changes**: purely additive
