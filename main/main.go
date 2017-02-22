package main

import (
	"hw"
	"config"
	//"fmt"
	//"time"
)

type HelloMsg struct {
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
			elev.ID
			
		}
	}
}

func main() {
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	elevTx := make(chan ElevatorInfo)
	elevRx := make(chan ElevatorInfo)

	go bcast.Transmitter(16569, helloTx, elevTx)
	go bcast.Receiver(16569, helloRx, elevRx)

	go f()

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

	Elev_init()
	Elev_set_motor_direction(DIR_UP)

	while(1) {
		if (Elev_get_floor_sensor_signal() == N_FLOORS - 1) {
			Elev_set_motor_direction(DIR_DOWN)
		} else if (Elev_get_floor_sensor_signal() == 0) {
			Elev_set_motor_direction(DIR_UP)
		}
		if(Elev_get_stop_signal()) {
			Elev_set_motor_direction(DIR_STOP)
			return 0
		}
	}
}