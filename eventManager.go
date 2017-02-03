package eventManager

import (
	"./config"
)

type Channels struct {
	NewOrder chan bool //Event
	ReachedFloor chan int //Event
	DoorClosing chan bool //Enten newOrder eller Idle

}

var floor int
var dir int

func eventManager(ch Channels) {
	for {
		select {
		case <-ch.NewOrder:
			handleNewOrder(ch)
		case floor := <-ch.ReachedFloor:
			handleReachedFloor(ch, floor)
		case <- DoorClosing:
			handleDoorClosing(ch)
		}
	}
}

func handleNewOrder(ch Channels) { //Til stud.ass, hvor tar vi inn info om pakker?
	switch state {
	case Idle:
		/*  Is the order in current floor?
				Open door
				state = OpenDoor
			Set motor direction
			state = moving

		*/
	}
	case Moving:
		//Ignore
	case OpenDoor:
		/*Set door lamp
				Start counter
				when counter out
					close Door
					turn off door lamp
					Delete order from queue
					*/

}

func handleReachedFloor(ch Channels) { //Til stud.ass - NewFloor hvordan?
	floor = newFloor
	Elev_set_floor_indicator(floor)
	switch state {
	case Moving:
		/*if correctFloor:
				stop
				Open door
				state = OpenDoor*/
	}
}