# Sense HAT I2C Reference

## I2C Addresses

| Sensor | Chip | I2C Address | Datasheet |
|---|---|---|---|
| Pressure/temperature | STMicro LPS25HB | `0x5C` | [LPS25HB](https://www.st.com/resource/en/datasheet/lps25hb.pdf) |
| Humidity/temperature | STMicro HTS221 | `0x5F` | [HTS221](https://www.st.com/resource/en/datasheet/hts221.pdf) |
| IMU (accel/gyro) | STMicro LSM9DS1 | `0x6A` | [LSM9DS1](https://www.st.com/resource/en/datasheet/lsm9ds1.pdf) |
| IMU (magnetometer) | STMicro LSM9DS1 | `0x1C` | (same as above) |
| Colour/brightness | TCS3400 | `0x29` or `0x39` | [TCS3400](https://ams.com/documents/20143/36005/TCS3400_DS000524_5-00.pdf) |
| LED matrix + joystick | ATTINY88 | `0x46` | - |

The I2C bus on the Raspberry Pi is typically available at `/dev/i2c-1`.

## Key Registers

### LPS25HB (Pressure)

- `0x28`, `0x29`, `0x2A` — 24-bit pressure output (XL, L, H)
- `0x2B`, `0x2C` — 16-bit temperature output (L, H)
- `0x20` — CTRL_REG1 (power on, data rate)
- `0x0F` — WHO_AM_I (should return `0xBD`)

### HTS221 (Humidity)

- `0x28`, `0x29` — 16-bit humidity output (L, H)
- `0x2A`, `0x2B` — 16-bit temperature output (L, H)
- `0x20` — CTRL_REG1 (power on, data rate)
- `0x0F` — WHO_AM_I (should return `0xBC`)
- `0x30`-`0x3F` — Calibration registers (required for converting raw values)

### LSM9DS1 (IMU)

Accelerometer/Gyroscope (at `0x6A`):
- `0x28`-`0x2D` — Gyroscope X, Y, Z output (16-bit each)
- `0x28`-`0x2D` — Accelerometer X, Y, Z output (16-bit each, different register bank)
- `0x0F` — WHO_AM_I (should return `0x68`)

Magnetometer (at `0x1C`):
- `0x28`-`0x2D` — Magnetometer X, Y, Z output (16-bit each)
- `0x0F` — WHO_AM_I (should return `0x3D`)

### TCS3400 (Colour)

- `0x94`-`0x95` — Clear channel (16-bit)
- `0x96`-`0x97` — Red channel (16-bit)
- `0x98`-`0x99` — Green channel (16-bit)
- `0x9A`-`0x9B` — Blue channel (16-bit)
- `0x80` — ENABLE register
- `0x92` — ID register

## Communication Architecture

The Sense HAT has two communication patterns:

1. **LED matrix and joystick** are mediated by an ATTINY88 microcontroller (`0x46`).
   The ATTINY exposes 256 registers: addresses 0-191 for the framebuffer
   (64 pixels x 3 bytes RGB), and `0xF2` for joystick state.
   BobbleHat currently accesses these through Linux kernel drivers
   (`/dev/fb*` for the display, `/dev/input/event*` for the joystick).

2. **Sensors** (pressure, humidity, IMU, colour) communicate directly with
   the Raspberry Pi over I2C. Adding Go support means reading/writing
   registers on `/dev/i2c-1` at the addresses listed above.

## I2C from Go

Options for I2C access in Go:

- **`periph.io/x/conn/v3/i2c`** — part of the periph.io hardware abstraction library, actively maintained
- **`golang.org/x/exp/io/i2c`** — experimental, minimal, may be sufficient for simple register access
- **Direct syscalls** — open `/dev/i2c-1`, use `ioctl` to set slave address, then `read`/`write` for register access (similar to how `sense/stick` already works with evdev)

## References

- [Official Python sense-hat library](https://github.com/astro-pi/python-sense-hat)
- [RTIMULib (C++/Python sensor abstraction)](https://github.com/RPi-Distro/RTIMULib)
- [Sense HAT I2C protocol documentation](https://github.com/underground-software/sensehat/blob/master/i2c_documentation.md)
- [Sense HAT product brief (PDF)](https://pip-assets.raspberrypi.com/categories/676-raspberry-pi-sense-hat/documents/RP-008384-DS-1-sense-hat-product-brief.pdf)
