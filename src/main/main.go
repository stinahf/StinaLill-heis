package main

import (
	"../config"
	"../eventManager"
	"../hw"
	"../queue"
	//"fmt"
	"time"
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

	ch := eventManager.Channels{
		NewOrder: make(chan bool),
		ReachedFloor: make(chan int),
		Message: make(chan config.Message),
		MotorDir: make(chan int),
		DoorLamp: make(chan bool),
		//DoorTimerReset: make(chan bool),
		//DoorTimeout: make(chan bool),
	}

	hw.Init()
	eventManager.Init(ch)
	hw.SetMotorDirection(config.DIR_DOWN)
	queue.Init(ch.NewOrder)

	go manageEvents(ch)

	time.Sleep(time.Second * 300)

}

func manageEvents(ch eventManager.Channels) {
	buttonPress := eventManager.PollButtons()
	floorHIT := eventManager.PollFloors()
	for {
		select {
		case button := <-buttonPress:
			switch button.Button {
			case config.BUTTON_UP, config.BUTTON_DOWN:
				queue.AddLocalOrder(button.Floor, button.Button, 1000)
			case config.BUTTON_INTERNAL:
				queue.AddLocalOrder(button.Floor, button.Button, 1000)
			}
		case floor := <-floorHIT:
			ch.ReachedFloor <- floor
		case dir := <-ch.MotorDir:
			hw.SetMotorDirection(dir)
		case value := <-ch.DoorLamp:
			hw.SetDoorOpenLamp(value)
		}
	}
}
