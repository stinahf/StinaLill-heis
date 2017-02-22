package eventManager

import (
	"./config"
	"./queue"
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

func EventManagerInit(ch Channels) {
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
	Elev_set_button_lamp(floor, button, true)
	switch state {
	case Idle:
		if shouldStop(def.DIR_STOP, floor){
			openDoor()
			Local_queue.RemoveOrder(floor, ch.Message) //Hvordan slette noe fra køen?
			handleDoorClosing(floor, dir)
		}
		else {
			dir = chooseMotorDirection(floor, dir)
			Elev_set_motor_direction(dir)
			state = Moving
		}

	}
	case Moving:
		//Ignore
	case OpenDoor:
		

}

func handleReachedFloor(ch Channels) {
	floor = newFloor
	Elev_set_floor_indicator(floor)
	switch state {
	case Idle:
		//Ignore
	case Moving:
		if shouldStop(dir, floor) {
			Elev_set_motor_direction(def.DIR_STOP)
			openDoor()
			Local_queue.RemoveOrder(floor, ch.Message)
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
		Elev_set_motor_direction(dir)
		state = Moving
	}
}



func openDoor() {
	Elev_set_door_open_lamp(true)
	timer := time.NewTimer(time.Second * 3)
	<- timer.C
	Elev_set_door_open_lamp(false)
}