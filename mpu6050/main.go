package main

import (
	"time"
	"fmt"

	"github.com/hybridgroup/gobot"
	//"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/hybridgroup/gobot/platforms/i2c"
	//"github.com/hybridgroup/gobot/platforms/firmata"
)

//dt 循环间隔,用于积分
var dt float64 = 0.01
var angle float64
var weight float64 = 0.98

func main() {
	gbot := gobot.NewGobot()

	raspiAdaptor := raspi.NewRaspiAdaptor("/dev/i2c-1")
	mpu6050 := i2c.NewMPU6050Driver(raspiAdaptor, "mpu6050")

	work := func() {
		gobot.Every(time.Millisecond*time.Duration(dt*1000), func() {
			//fmt.Println("Accelerometer", mpu6050.Accelerometer,"Gyroscope", mpu6050.Gyroscope,"Temperature", mpu6050.Temperature)
			angle = weight*(angle + float64(mpu6050.Gyroscope.Y)*dt)+(1-weight)*float64(mpu6050.Accelerometer.X)
			fmt.Println(int(angle/100))
		})
	}

	robot := gobot.NewRobot("mpu6050Bot",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{mpu6050},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}