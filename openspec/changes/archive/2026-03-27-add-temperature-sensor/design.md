## Context

Bobblehat currently supports the Sense HAT LED matrix (via framebuffer) and joystick (via evdev). The Sense HAT also has weather sensors accessible over I2C, but no Go support exists for them yet. The LPS25HB pressure/temperature sensor at I2C address `0x5C` is the simplest entry point — temperature requires only a register read and linear conversion with no calibration.

The existing `sense/stick` package demonstrates direct syscall-based device access. For sensors, we use `periph.io/x/conn/v3/i2c` interfaces instead, which provides testability and follows established Go hardware library conventions.

## Goals / Non-Goals

**Goals:**
- Provide a `sense/sensors` package with a `Temperature()` method returning °C as `float64`
- Use `periph.io/x/conn/v3/i2c` for I2C bus abstraction
- Establish a pattern that pressure, humidity, and IMU sensors can follow later
- Include a working example program

**Non-Goals:**
- Pressure reading (future addition to same package)
- Humidity reading via HTS221 (future addition)
- IMU / motion sensor support
- Auto-discovery of I2C bus (application provides the bus)
- Continuous/streaming reads — this is a poll-based API

## Decisions

### 1. Use `periph.io/x/conn/v3/i2c` for I2C access

**Choice**: Depend on periph.io interfaces rather than raw syscalls.

**Alternatives considered**:
- **Raw syscalls** (like `stick.go`): zero deps, but no testability without hand-rolling a bus interface. We'd end up reinventing what conn/v3 provides.
- **`golang.org/x/exp/io/i2c`**: experimental, unclear future.

**Rationale**: `conn/v3` adds near-zero weight (interfaces + stdlib only). Applications import `host/v3` separately to get real drivers. This also makes unit testing possible with a mock bus.

### 2. Package name: `sense/sensors`

**Choice**: Single `sensors` package, organized by measurement type (like Python's `sense-hat`).

**Alternatives considered**:
- **`sense/weather`**: too narrow if IMU support is added later
- **Per-chip packages** (`sense/lps25hb`): exposes hardware details users don't care about

**Rationale**: Users want `Temperature()`, not `LPS25HB.ReadRegister()`. The package can grow to include `Pressure()`, `Humidity()`, etc.

### 3. API shape: Open/close with poll-based reads

```go
dev, err := sensors.Open(bus)   // bus is an i2c.Bus from periph.io
defer dev.Close()
temp, err := dev.Temperature()  // returns float64 in °C
```

**Rationale**: Matches the `stick.Open()` pattern. Poll-based (not channel-based) because sensor reads are synchronous I2C transactions — no background goroutine needed.

### 4. LPS25HB temperature conversion

The sensor outputs a 16-bit signed integer. Conversion: `float64(int16(raw)) / 480.0 + 42.5`.

No calibration registers required — this is defined by the datasheet.

## Risks / Trade-offs

- **[Hardware unavailable during development]** → Use mock `i2c.Bus` in tests; verify on real hardware separately
- **[LPS25HB temperature is influenced by board heat]** → Known Sense HAT limitation; not something we can fix in software. Document it.
- **[periph.io dependency could become unmaintained]** → Low risk; conn/v3 is stable and interface-only. Could be replaced without API changes if needed.
