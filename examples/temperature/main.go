package main

import (
	"fmt"
	"log"

	"github.com/tonvanbart/bobblehat/sense/sensors"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func main() {
	// Initialize periph.io host drivers.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open the default I2C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Open the Sense HAT sensors.
	dev, err := sensors.Open(bus)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Read and print the temperature.
	temp, err := dev.Temperature()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Temperature: %.1f°C\n", temp)
}
