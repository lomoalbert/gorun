package mpuserial

import (
	"github.com/tarm/serial"
	"math"
	"log"
)

var xrolloffset float64 = 2.6
var yacceoffset float64 = 0.4

func read(port *serial.Port, i int) []byte {
	b := make([]byte, i)
	b1 := make([]byte, 1)
	var sum int
	for t := 0; t < i; t++ {
		n, _ := port.Read(b1)
		b[t] = b1[0]
		sum += n
	}
	return b
}

func align(port *serial.Port) {
	log.Println("serial port align byte")
	b := make([]byte, 1)
	for {
		port.Read(b)
		if b[0] == 0x55 {
			read(port, 10)
			return
		}
	}
}

func getyacce(frame []byte) (float64, error) {
	yacce := float64(int16(frame[5]) << 8 | int16(frame[4])) * 16 * 9.8 / 32768
	return yacce, nil
}

func getxroll(frame []byte) (float64, error) {
	xroll := float64(int64(frame[3]) << 8 | int64(frame[2])) * 180 / 32768
	return xroll, nil
}

func getxspeed(frame []byte) (float64, error) {
	xspeed := float64(int16(frame[3]) << 8 | int16(frame[2])) * 2000 / 32768
	return xspeed, nil
}

type MPU struct {
	Conf   *serial.Config
	Port   *serial.Port
	XRoll  float64
	Yacce  float64
	Xspeed float64
	PWM    float64
	K1     float64  //角度权重
	K2     float64  //加速度权重
	K3     float64  //速度权重
}

func NewMPU() (*MPU, error) {
	var err error
	mpu := new(MPU)
	mpu.Conf = &serial.Config{Name: "/dev/ttyAMA0",
		Baud: 115200,
		Size: serial.DefaultSize,
		Parity:serial.ParityNone,
		StopBits:serial.Stop1,
	}
	mpu.Port, err = serial.OpenPort(mpu.Conf)
	if err != nil {
		return nil, err
	}
	mpu.PWM = 0.0
	return mpu, nil
}

func (mpu *MPU)Start() {
	log.Println("Starting")
	for {
		buf := read(mpu.Port, 11)
		if buf[0]!=0x55{
			log.Println(buf)
			align(mpu.Port)
			continue
		}
		if buf[1] == 0x53 {
			xroll, err := getxroll(buf)
			xroll = xroll - xrolloffset
			if err != nil {
				continue
			}
			if xroll > 180.0 {
				mpu.XRoll = xroll - 360.0
			} else {
				mpu.XRoll = xroll
			}
		} else if buf[1] == 0x51 {
			yacce, err := getyacce(buf)
			yacce = yacce - yacceoffset
			if err != nil {
				continue
			}
			mpu.Yacce = yacce
		}else if buf[1] == 0x52 {
			xspeed,err := getxspeed(buf)
			if err != nil {
				continue
			}
			mpu.Xspeed = xspeed
		}

		pwm :=  mpu.XRoll*mpu.K1 + mpu.Yacce*mpu.K2 + mpu.Xspeed *mpu.K3
		//go log.Println(mpu.K1,mpu.K2,mpu.XRoll,mpu.Yacce,mpu.Yspped,pwm)
		var dir float64 = 1
		if pwm <0 {
			dir = -1
		}
		pwm =  math.Min(math.Abs(pwm), 1) * dir
		mpu.PWM = pwm
	}
}