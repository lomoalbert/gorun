package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

//dt 循环间隔,用于积分
var dt float64 = 0.001
var angle float64
var weight float64 = 0.995
var offset float64 = 110

var raspiAdaptor = raspi.NewRaspiAdaptor("raspi")
var pin16 = NewStepPin(gpio.NewDirectPinDriver(raspiAdaptor,"pins","16")) //A相
var pin12 = NewStepPin(gpio.NewDirectPinDriver(raspiAdaptor,"pins","16")) //A相

var pin29 = NewPWMPin(gpio.NewDirectPinDriver(raspiAdaptor,"pin29","29"),pin16)
var pin31 = NewPWMPin(gpio.NewDirectPinDriver(raspiAdaptor,"pin31","31"),pin16)
var pin35 = NewPWMPin(gpio.NewDirectPinDriver(raspiAdaptor,"pin35","35"),pin12)
var pin37 = NewPWMPin(gpio.NewDirectPinDriver(raspiAdaptor,"pin37","37"),pin12)



var frontlist = []*PWMPin{pin29,  pin35}
var backlist = []*PWMPin{pin31,pin37}

func speedpin(pinlist []*PWMPin,speed float64){
	for _,pin := range pinlist{
		pin.Power(speed)
	}
}

func motor(direction int,speed float64) {
	switch  {
	case direction > 0:
		speedpin(frontlist,speed)
		speedpin(backlist,0)
	case direction < 0:
		speedpin(backlist,speed)
		speedpin(frontlist,0)
	case direction == 0:
		speedpin(frontlist,0)
		speedpin(backlist,0)
	}
}

func slow(speed float64)float64{
	return speed/150/100000
}

func policy(angle float64){
	var boundary float64 = 150
	if angle > 4000 || angle < -4000{
		motor(0,0)
		return
	}
	switch  {
	case angle < boundary && angle > -boundary:
		motor(0,0)
	case angle > boundary:
		motor(-1,slow(angle-boundary))
	case angle < -boundary:
		motor(1,slow(boundary-angle))
	}

}

var i int
func Balance(acce,gyro float64)float64{
	if i%100==0{
		go println(weight*(angle + gyro*dt)+acce*(1-weight))
	}
	i+=1
	return weight*(angle + gyro*dt)+acce*(1-weight)
}

func main() {
	gbot := gobot.NewGobot()

	mpu6050 := i2c.NewMPU6050Driver(raspiAdaptor, "mpu6050",0)

	work := func() {
		gobot.Every(time.Millisecond*time.Duration(dt*1000), func() {
			angle = Balance(float64(mpu6050.Accelerometer.X)+offset,float64(mpu6050.Gyroscope.Y))
			policy(angle)
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