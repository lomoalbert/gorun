package main

import (
	"time"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"fmt"
)


var r = raspi.NewRaspiAdaptor("raspi")
var pin31 = gpio.NewLedDriver(r, "led", "31")
var pin33 = gpio.NewLedDriver(r, "led", "33")
var pin35 = gpio.NewLedDriver(r, "led", "35")
var pin37 = gpio.NewLedDriver(r, "led", "37")

func work() {
	pinlist := []*gpio.LedDriver{pin31, pin33, pin35, pin37}
	gobot.Every(time.Millisecond*100 * 4, func() {
		for n, pin := range pinlist {
			fmt.Println(n)
			pin.On()
			time.Sleep(time.Millisecond*100)
			pin.Off()
		}
	})
}

func main() {

	gbot := gobot.NewGobot()

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{pin31, pin33, pin35, pin37},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
