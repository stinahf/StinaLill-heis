package eventManager

import (
	"../config"
	"../hw"
	"../queue"
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
}

func GetFloorDirState() (int, int, int) {
	return floor, dir, state
}

func Run(ch Channels) {
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
	switch state {
	case config.Idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(dir, floor) {
			ch.DoorTimerReset <- true
			state = config.DoorOpen
			queue.RemoveOrder(floor)
			ch.DoorLamp <- true
		} else {
			ch.MotorDir <- dir
			state = config.Moving
		}
	case config.Moving:
		//Ignore
	case config.DoorOpen:
		if queue.ShouldStop(dir, floor) {
			ch.DoorTimerReset <- true
			queue.RemoveOrder(floor)
		}
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
			state = config.DoorOpen
			queue.RemoveOrder(floor)
			ch.DoorLamp <- true
			ch.DoorTimerReset <- true
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
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-resetTimer:
			timer.Reset(3 * time.Second)
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
