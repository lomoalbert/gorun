package main

import (
	"time"
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

//dt 循环间隔,用于积分
var dt float64 = 0.01
var angle float64
var weight float64 = 0.97


var raspiAdaptor = raspi.NewRaspiAdaptor("raspi")
var motorcontroller = gpio.NewMotorDriver(raspiAdaptor, "motor", "3")
var pin29 = gpio.NewLedDriver(raspiAdaptor, "led", "29")
var pin31 = gpio.NewLedDriver(raspiAdaptor, "led", "31")
var pin35 = gpio.NewLedDriver(raspiAdaptor, "led", "35")
var pin37 = gpio.NewLedDriver(raspiAdaptor, "led", "37")


var frontlist = []*gpio.LedDriver{pin29,  pin35,}
var backlist = []*gpio.LedDriver{pin31,pin37}

func onpin(pinlist []*gpio.LedDriver){
	for _,pin := range pinlist{
		pin.On()
	}
}

func offpin(pinlist []*gpio.LedDriver){
	for _,pin := range pinlist{
		pin.Off()
	}
}

func motor(direction int) {
	frontlist := []*gpio.LedDriver{pin29,  pin35,}
	backlist := []*gpio.LedDriver{pin31,pin37}
	switch  {
	case direction > 0:
		offpin(backlist)
		onpin(frontlist)
	case direction < 0:
		offpin(frontlist)
		onpin(backlist)
	case direction == 0:
		offpin(frontlist)
		offpin(backlist)
	}
}

func policy(angle float64){
	switch  {
	case angle < 800 && angle > -800:
		motor(0)
	case angle > 1000:
		motor(-1)
	case angle < -1000:
		motor(1)
	}

}

func main() {
	gbot := gobot.NewGobot()

	mpu6050 := i2c.NewMPU6050Driver(raspiAdaptor, "mpu6050")

	work := func() {
		gobot.Every(time.Millisecond*time.Duration(dt*1000), func() {
			//fmt.Println("Accelerometer", mpu6050.Accelerometer,"Gyroscope", mpu6050.Gyroscope,"Temperature", mpu6050.Temperature)
			angle = weight*(angle + float64(mpu6050.Gyroscope.Y)*dt)+(1-weight)*float64(mpu6050.Accelerometer.X)
			policy(angle)
		})
	}
	robot := gobot.NewRobot("mpu6050Bot",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{mpu6050},
		work,
	)
	fmt.Println(robot)
	gbot.AddRobot(robot)
	gbot.Start()
}