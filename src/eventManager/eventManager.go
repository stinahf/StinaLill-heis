package eventManager

import (
	def "config"
	"queue"
	"elev"
	"time"
)

type Channels struct {
	NewOrder chan bool //Event
	ReachedFloor chan int //Event
	DoorClosing chan bool //Enten newOrder eller Idle
	Message chan def.Message
}

var floor int
var dir int

func Init(ch Channels) {
	state = Idle
	dir = def.DIR_STOP
	floor = 0

	go eventManager(ch)

	
}

func eventManager(ch Channels) {
	for {
		select {
		case <-ch.NewOrder:
			handleNewOrder(ch)
		case floor := <-ch.ReachedFloor:
			handleReachedFloor(ch, floor)
		}
	}
}

func handleNewOrder(order OrderInfo) {
	hw.SetButtonLamp(floor, button, true)
	switch state {
	case Idle:
		if shouldStop(def.DIR_STOP, floor){
			openDoor()
			queue.RemoveOrder(floor, ch.Message) //Fikset, sjekk at funker
			handleDoorClosing(floor, dir)
		}
		else {
			dir = chooseMotorDirection(floor, dir)
			hw.SetMotorDirection(dir)
			state = Moving
		}

	}
	case Moving:
		//Ignore
	case OpenDoor:
		

}

func handleReachedFloor(ch Channels) {
	floor = newFloor
	hw.SetFloorIndicator(floor)
	switch state {
	case Idle:
		//Ignore
	case Moving:
		if shouldStop(dir, floor) {
			hw.SetMotorDirection(def.DIR_STOP)
			openDoor()
			queue.RemoveOrder(floor, ch.Message)
			handleDoorClosing(floor, dir)
		}
	case OpenDoor:
		//Ignore			
	}
}

func handleDoorClosing(floor int, dir int) {
	if ChooseMotorDirection(floor,dir) == def.DIR_STOP {
		state = Idle
	}
	else {
		dir = ChooseMotorDirection(floor, dir)
		hw.SetMotorDirection(dir)
		state = Moving
	}
}



func openDoor() {
	hw.SetDoorOpenLamp(true)
	timer := time.NewTimer(time.Second * 3)
	<- timer.C
	hw.SetDoorOpenLamp(false)
}