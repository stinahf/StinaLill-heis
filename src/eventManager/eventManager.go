package eventManager

import (
	"../config"
	"../queue"
	"../hw"
	"time"
)

type Channels struct {
	NewOrder chan bool //Event
	ReachedFloor chan int //Event
	DoorClosing chan bool //Enten newOrder eller Idle
	Message chan config.Message
}

var floor int
var dir int
var state int

func Init(ch Channels) {
	state = config.Idle
	dir = config.DIR_STOP
	floor = 0

	go eventManager(ch)

	
}

func eventManager(ch Channels) {
	for {
		select {
		case <-ch.NewOrder:
			handleNewOrder(ch)
		case floor := <-ch.ReachedFloor:
			handleReachedFloor(ch)
		}
	}
}

func handleNewOrder(ch Channels) {
	floor = queue.GetFloorFromQueue()
	button = queue.GetButtonFromQueue()
	hw.SetButtonLamp(floor, button, true)
	switch state {
	case Idle:
		if queue.ShouldStop(config.DIR_STOP, floor){
			openDoor()
			queue.RemoveOrder(floor, ch.Message) //Fikset, sjekk at funker
			floor = queue.GetFloorFromQueue()
			handleDoorClosing(floor, config.DIR_STOP)
		} else {
			dir = chooseMotorDirection(floor, dir)
			hw.SetMotorDirection(dir)
			state = Moving
		}
	case Moving:
		//Ignore
	case OpenDoor:
	}

}

func handleReachedFloor(ch Channels) {
	floor = newFloor
	hw.SetFloorIndicator(floor)
	switch state {
	case Idle:
		//Ignore
	case Moving:
		if shouldStop(dir, floor) {
			hw.SetMotorDirection(config.DIR_STOP)
			openDoor()
			queue.RemoveOrder(floor, ch.Message)
			//Få tak i direction from infoPackage her
			handleDoorClosing(floor, dir)
		}
	case OpenDoor:
		//Ignore			
	}
}

func handleDoorClosing(floor int, dir int) {
	if queue.ChooseMotorDirection(floor,dir) == config.DIR_STOP {
		state = Idle
	} else {
		dir = queue.ChooseMotorDirection(floor, dir)
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
