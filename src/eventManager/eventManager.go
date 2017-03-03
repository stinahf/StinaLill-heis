package eventManager

import (
	"../config"
	"../hw"
	"../queue"
	"time"
)

type Channels struct {
	NewOrder     chan bool //Event
	ReachedFloor chan int  //Event
	DoorClosing  chan bool //Enten newOrder eller Idle
	Message      chan config.Message
}

var floor int
var newFloor int
var button int
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
		case <-ch.ReachedFloor:
			handleReachedFloor(ch)
		}
	}
}

func handleNewOrder(ch Channels) {
	/*floor = queue.GetFloorFromQueue()
	button = queue.GetButtonFromQueue()*/
	hw.SetButtonLamp(floor, button, true)
	switch state {
	case config.Idle:
		if queue.ActuallyShouldStop(config.DIR_STOP, floor) {
			openDoor()
			queue.RemoveOrder(floor, ch.Message) //Fikset, sjekk at funker
			floor = queue.GetFloorFromQueue()
			handleDoorClosing(floor, config.DIR_STOP)
		} else {
			dir = queue.ActuallyChooseDirection(floor, dir)
			hw.SetMotorDirection(dir)

			state = config.Moving
		}
	case config.Moving:
		//Ignore
	case config.OpenDoor:
	}

}

func handleReachedFloor(ch Channels) {
	/*floor = hw.GetFloorSensorSignal()
	hw.SetFloorIndicator(floor)*/
	switch state {
	case config.Idle:
		//Ignore
	case config.Moving:
		if queue.ActuallyShouldStop(dir, floor) {
			hw.SetMotorDirection(config.DIR_STOP)
			openDoor()
			queue.RemoveOrder(floor, ch.Message)
			//Få tak i direction from infoPackage her
			handleDoorClosing(floor, dir)
		}
	case config.OpenDoor:
		//Ignore
	}
}

func handleDoorClosing(floor int, dir int) {
	if queue.ActuallyChooseDirection(floor, dir) == config.DIR_STOP {
		state = config.Idle
	} else {
		dir = queue.ActuallyChooseDirection(floor, dir)
		hw.SetMotorDirection(dir)
		state = config.Moving
	}
}

func openDoor() {
	hw.SetDoorOpenLamp(true)
	timer := time.NewTimer(time.Second * 3)
	<-timer.C
	hw.SetDoorOpenLamp(false)
}

func pollFloors() {

	oldFloor := hw.GetFloorSensorSignal()
	go func() {
		for {
			newFloor := hw.GetFloorSensorSignal()
			if newFloor != oldFloor && newFloor != -1 {
				ReachedFloor <- newFloor
			}

			oldFloor = newFloor
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func pollButtons() <-chan config.OrderInfo {
	buttonPress := make(chan config.OrderInfo)
	go func() {
		for {
			for floor := 0; floor < config.N_FLOORS; floor++ {
				for button := 0; button < config.N_BUTTONS; button++ {
					if hw.GetButtonSignal(floor, button) {
						buttonPress <- config.OrderInfo{Button: button, Floor: floor}
					}
				}
			}
		}
	}()
	return buttonPress		
}
