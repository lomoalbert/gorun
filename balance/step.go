package main

import "time"

type StepPiner interface {
	Name() string
	DigitalRead() (val int, err error)
}

type StepPin struct {
	pin StepPiner
	interval time.Duration
	status int
}

func NewStepPin(sp StepPiner,v ...time.Duration)*StepPin{
	steppin := &StepPin{
		pin: sp,
		interval: 0,
		status:0,
	}
	if len(v)>0{
		steppin.interval = v[0]
	}
	return steppin
}

func (SP StepPin)Count(step int){
	SP.status ,_ = SP.pin.DigitalRead()
	step *= 2
	for step>0 {
		tick := time.Tick(SP.interval)
		status ,err := SP.pin.DigitalRead()
		if err == nil{
			if status != SP.status{
				SP.status = status
				step--
			}
		}
		if tick != nil{
			<- tick
		}
	}
}

