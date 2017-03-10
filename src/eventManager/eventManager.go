package eventManager

import (
	"../config"
	"../hw"
	"../queue"
	"fmt"
	"time"
)

type Channels struct {
	NewOrder       chan bool
	DoorLamp       chan bool
	ReachedFloor   chan int
	MotorDir       chan int
	DoorTimeout    chan bool
	DoorTimerReset chan bool
}

var floor int
var dir int
var state int

func Init() {

	state = config.Idle
	dir = config.DIR_STOP
	floor = 0

	/*ch.DoorTimeout = make(chan bool)
	ch.DoorTimerReset = make(chan bool)*/

}

func GetFloorDirState() (int, int, int) {
	return floor, dir, state
}

func EventManager(ch Channels) {
	for {
		select {
		case <-ch.NewOrder:
			handleNewOrder(ch)
		case floor := <-ch.ReachedFloor:
			handleReachedFloor(ch, floor)
		case <-ch.DoorTimeout:
			handleDoorClosing(ch)
		}
	}
}

func handleNewOrder(ch Channels) {

	fmt.Println("Jeg skal utføre en ordre!")

	switch state {
	case config.Idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(dir, floor) {
			ch.DoorTimerReset <- true
			queue.RemoveOrder(floor)
			state = config.DoorOpen
			ch.DoorLamp <- true
		} else {
			ch.MotorDir <- dir
			state = config.Moving
		}
	case config.Moving:
		//Ignore
	}

}

func handleReachedFloor(ch Channels, newFloor int) {
	floor = newFloor
	switch state {
	case config.Idle:
		//Ignore
	case config.Moving:
		if queue.ShouldStop(dir, floor) {
			dir = config.DIR_STOP
			ch.MotorDir <- dir
			queue.RemoveOrder(floor)
			state = config.DoorOpen
			ch.DoorTimerReset <- true
			ch.DoorLamp <- true
		}
	}
}

func handleDoorClosing(ch Channels) {
	ch.DoorLamp <- false

	if queue.ChooseDirection(floor, dir) == config.DIR_STOP {
		dir = config.DIR_STOP
		ch.MotorDir <- dir
		state = config.Idle
	} else {
		dir = queue.ChooseDirection(floor, dir)
		ch.MotorDir <- dir
		state = config.Moving
	}
}

func OpenDoor(doorTimeout chan<- bool, resetTimer <-chan bool) {
	hw.SetDoorOpenLamp(true)
	timer := time.NewTimer(0)
	timer.Stop()
	hw.SetDoorOpenLamp(false)
	for {
		select {
		case <-resetTimer:
			hw.SetDoorOpenLamp(true)
			timer.Reset(3 * time.Second)
			hw.SetDoorOpenLamp(false)
		case <-timer.C:
			timer.Stop()
			doorTimeout <- true
		}
	}

}

func PollFloors(temp chan int) {
	oldFloor := hw.GetFloorSensorSignal()
	for {
		newFloor := hw.GetFloorSensorSignal()
		if newFloor != oldFloor && newFloor != -1 {
			hw.SetFloorIndicator(newFloor)
			temp <- newFloor
		}
		oldFloor = newFloor
		time.Sleep(time.Millisecond * 100)
	}
}

func PollButtons(temp chan config.OrderInfo) {
	var pressed [config.N_FLOORS][config.N_BUTTONS]bool
	for {
		for floor := 0; floor < config.N_FLOORS; floor++ {
			for button := 0; button < config.N_BUTTONS; button++ {
				if (floor == 0 && button == config.BUTTON_DOWN) || (floor == config.N_FLOORS-1 && button == config.BUTTON_UP) {
					continue
				}
				if hw.GetButtonSignal(button, floor) {
					if !pressed[floor][button] {
						pressed[floor][button] = true
						temp <- config.OrderInfo{Button: button, Floor: floor}
						hw.SetButtonLamp(button, floor, true)
					}
				} else {
					pressed[floor][button] = false
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}
