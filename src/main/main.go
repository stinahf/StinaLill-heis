package main

import (
	def "config"
	"hw"
	//"fmt"
	//"time"
)

/*type HelloMsg struct {
	Message string
	Iter int
}

func f() {
	for {
		select {
		case hello := <-helloRx:
			// do some code
		case elev := <-elevRx:
			// do some other code
			//elev.ID

		}
	}
}*/

func main() {
	/*
		helloTx := make(chan HelloMsg)
		helloRx := make(chan HelloMsg)
		elevTx := make(chan ElevatorInfo)
		elevRx := make(chan ElevatorInfo)

		go bcast.Transmitter(16569, helloTx, elevTx)
		go bcast.Receiver(16569, helloRx, elevRx)

		go f()
	*/
	//go func() {
	/*	helloMsg := HelloMsg{"Hello from us", 0}
			for {
				helloMsg.Iter++
				helloTx <- helloMsg
				time.Sleep(1 * time.Second)
			}
		}//()

		fmt.Println("Started")
		for {
			select {
			case a := <-helloRx:
				fmt.Printf("Received: %#v\n", a)
			}
		}*/

	hw.Init() //Gir ut feil meldinger: You are already at the bottom, Invalid button 3 to ganger, tyder på at heisen tror knapper blir trykket inn når det ikke blir det
	hw.SetMotorDirection(def.DIR_UP)

	for { //Go har ikke while loops
		if hw.GetFloorSensorSignal() == def.N_FLOORS-1 {
			hw.SetMotorDirection(def.DIR_DOWN)
		} else if hw.GetFloorSensorSignal() == 0 {
			hw.SetMotorDirection(def.DIR_UP)
		}

		if hw.GetStopSignal() {
			hw.SetMotorDirection(def.DIR_STOP)
			//return 0
		}
	}
}
