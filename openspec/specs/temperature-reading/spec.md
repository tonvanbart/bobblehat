### Requirement: Open sensor device
The `sensors` package SHALL provide an `Open(bus i2c.Bus)` function that initializes the LPS25HB sensor and returns a device handle. The function MUST power on the sensor by writing to its control register.

#### Scenario: Successful open
- **WHEN** `sensors.Open(bus)` is called with a valid I2C bus
- **THEN** it returns a `*Device` and nil error, and the LPS25HB sensor is powered on with continuous output enabled

#### Scenario: WHO_AM_I verification
- **WHEN** `sensors.Open(bus)` is called
- **THEN** it MUST read the WHO_AM_I register (`0x0F`) and verify it returns `0xBD`, returning an error if the sensor is not detected

### Requirement: Read temperature
The `Device` SHALL provide a `Temperature()` method that returns the current temperature in degrees Celsius as a `float64`.

#### Scenario: Successful temperature read
- **WHEN** `Temperature()` is called on an open device
- **THEN** it reads the 16-bit temperature output registers (`0x2B`, `0x2C`) from the LPS25HB and returns the converted value in °C

#### Scenario: Temperature conversion
- **WHEN** the raw 16-bit sensor value is read
- **THEN** it SHALL be converted using the formula: `float64(int16(raw)) / 480.0 + 42.5`

### Requirement: Close device
The `Device` SHALL provide a `Close()` method for cleanup.

#### Scenario: Close after use
- **WHEN** `Close()` is called on an open device
- **THEN** the device handle is released and subsequent calls to `Temperature()` return an error

### Requirement: Example program
An example program SHALL be provided at `examples/temperature/main.go` that demonstrates reading temperature from the Sense HAT.

#### Scenario: Example reads and prints temperature
- **WHEN** the example program is run on a Raspberry Pi with a Sense HAT
- **THEN** it opens the I2C bus, reads the temperature, prints it to stdout, and exits cleanly
