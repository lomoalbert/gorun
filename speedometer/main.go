package speedometer

import (
	"fmt"
	"github.com/davecheney/gpio"
	"os"
	"os/signal"
	"log"
)

type Status struct {
	PinA  gpio.Pin
	PinB  gpio.Pin
	A  bool
	B  bool
	Count int64
}

func UpdateStatus(status *Status){
	status.A = status.PinA.Get()
	status.B = status.PinB.Get()
	if status.A && status.B{
		status.Count += 1
	}
}

func SpeedOmeter() *Status{
	var err error
	status := new(Status)
	log.Println(status)
	status.PinA, err = gpio.OpenPin(gpio.GPIO23, gpio.ModeInput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return nil
	}

	status.PinB, err = gpio.OpenPin(gpio.GPIO24, gpio.ModeInput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return nil
	}
	// clean up on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Println("Closing pin and terminating program.")
			status.PinA.Close()
			status.PinB.Close()
			os.Exit(0)
		}
	}()

	err = status.PinA.BeginWatch(gpio.EdgeBoth, func() {
		UpdateStatus(status)
	})
	if err != nil {
		fmt.Printf("Error beginwatch pin! %s\n", err)
		return nil
	}

	err = status.PinB.BeginWatch(gpio.EdgeBoth, func() {
		UpdateStatus(status)
	})
	if err != nil {
		fmt.Printf("Error beginwatch pin! %s\n", err)
		return nil
	}
	return status
}