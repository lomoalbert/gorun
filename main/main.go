package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"github.com/lomoalbert/gorun/mpuserial"
	"log"
	"math"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	pinlb = rpio.Pin(12)
	pinlf = rpio.Pin(13)
	pinrf = rpio.Pin(19)
	pinrb = rpio.Pin(26)
)



type Motor struct {
	pinlb rpio.Pin
	pinlf rpio.Pin
	pinrf rpio.Pin
	pinrb rpio.Pin
}

func NewMotor() *Motor {
	motor := new(Motor)
	motor.pinlf = pinlf
	motor.pinlb = pinlb
	motor.pinrf = pinrf
	motor.pinrb = pinrb
	return motor
}

func (m *Motor)init() {
	for _, pin := range []rpio.Pin{m.pinlb, m.pinlf, m.pinrf, m.pinrb} {
		pin.Output()
	}
}

func (m *Motor)run(mpu *mpuserial.MPU) {
	for {
		pwm := mpu.PWM
		if pwm > 0 {
			m.front()
		} else {
			m.back()
		}
		time.Sleep(time.Duration(math.Abs(pwm) * 1000) * time.Microsecond)
		m.stop()
		time.Sleep(time.Duration(1000.0 - math.Abs(pwm) * 1000) * time.Microsecond)
	}
}

func (m *Motor)stop() {
	for _, pin := range []rpio.Pin{m.pinlb, m.pinlf, m.pinrf, m.pinrb} {
		pin.Low()
	}
}

func (m *Motor)front() {
	for _, pin := range []rpio.Pin{pinlf, pinrf} {
		pin.High()
	}
}

func (m *Motor)back() {
	for _, pin := range []rpio.Pin{pinlb, pinrb} {
		pin.High()
	}
}

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()
	motor := NewMotor()
	motor.init()
	defer motor.stop()

	mpu, err := mpuserial.NewMPU()
	if err != nil {
		log.Println(err.Error())
		return
	}
	//mpu.K1 = 2.3/90
	//mpu.K2 = 0.3/9.8
	//mpu.K1 = 3.1/90
	//mpu.K2 = 0.3/9.8
	//mpu.K1 = 3.3/90
	//mpu.K2 = 0.5/9.8
	//mpu.K1 = 2.9/90
	//mpu.K2 = 0.33/9.8
	mpu.K1 = 1.0/15  //分母 即为PWM满偏时的角度，0-90
	mpu.K2 = 1.0/20 //分母 即为PWM满偏时的加速度，0-9.8
	mpu.K3 = 1.0/800
	go mpu.Start()
	// Set pin to output mode

	// Toggle pin 20 times
	//time.Sleep(time.Hour)
	motor.run(mpu)
	motor.stop()
}