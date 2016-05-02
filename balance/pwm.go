package main

type PWMer interface {
	Name() string
	DigitalWrite(level byte) (err error)
}

type PWMPin struct {
	pin   PWMer
	power float64
	steppin *StepPin
	counter int
}

func NewPWMPin(pin PWMer,steppin *StepPin) (*PWMPin) {
	return &PWMPin{
		pin:pin,
		power:0,
		steppin:steppin,
		counter: 0,
	}
}


//本电机建议驱动 PWM 频率是:10KHZ
func (pin *PWMPin)Run() {
	defer func(){
		if pin.counter == 0{
			pin.pin.DigitalWrite(0)
		}
	}()
	pin.counter++
	if pin.power != 0 {
		pin.pin.DigitalWrite(1)
		pin.steppin.Count(1)
	}
	pin.counter--
}

func (pin *PWMPin)Power(power float64) {
	switch  {
	case power < 0:
		power = 0
	case power > 1:
		power = 1
	}
	pin.power = power
	go pin.Run()
}