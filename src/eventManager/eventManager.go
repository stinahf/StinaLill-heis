package eventManager

import (
	"../config"
	"../hw"
	"../queue"
	"time"
	"fmt"
)

type Channels struct {
	NewOrder       chan bool
	DoorLamp       chan bool 
	Message        chan config.Message
	ReachedFloor   chan int
	MotorDir       chan int
	DoorTimeout    chan bool
	DoorTimerReset chan bool
}

floor := config.ElevatorInfo.CurrentFloor
dir   := config.ElevatorInfo.MotorDir
state := config.ElevatorInfo.State

func Init(ch Channels) {
	state = config.Idle
	dir = config.DIR_STOP
	floor = 0

	ch.DoorTimeout = make(chan bool)
	ch.DoorTimerReset = make(chan bool)

	go eventManager(ch)
	go OpenDoor(ch.DoorTimeout, ch.DoorTimerReset)

}

func Floor() int {
	return floor
}

func eventManager(ch Channels) {
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
		dir = queue.ActuallyChooseDirection(floor, dir)
		if queue.ActuallyShouldStop(dir, floor) {
			ch.DoorTimerReset <- true
			queue.RemoveOrder(floor)
			ch.DoorLamp <- true
		} else {
			ch.MotorDir <- dir
			state = config.Moving
		}
	case config.Moving:
		//Ignore
	case config.DoorClosing:
	}

}

func handleReachedFloor(ch Channels, newFloor int) {
	floor = newFloor
	switch state {
	case config.Idle:
		//Ignore
	case config.Moving:
		if queue.ActuallyShouldStop(dir, floor) {
			dir = config.DIR_STOP
			ch.MotorDir <- dir
			queue.RemoveOrder(floor)
			ch.DoorTimerReset <- true
			ch.DoorLamp <- true
		}
	case config.DoorClosing:
		//Ignore
	}
}

func handleDoorClosing(ch Channels) {
	ch.DoorLamp <- false

	if queue.ActuallyChooseDirection(floor, dir) == config.DIR_STOP {
		dir = config.DIR_STOP
		ch.MotorDir <- dir
		state = config.Idle
	} else {
		dir = queue.ActuallyChooseDirection(floor, dir)
		ch.MotorDir <- dir
		state = config.Moving
	}
}

func OpenDoor(doorTimeout chan <- bool, resetTimer <- chan bool) {
	hw.SetDoorOpenLamp(true)
	timer := time.NewTimer(0)
	timer.Stop()
	hw.SetDoorOpenLamp(false)

	for{
		select{
		case <- resetTimer:
			hw.SetDoorOpenLamp(true)
			timer.Reset(3*time.Second)
			hw.SetDoorOpenLamp(false)
		case <- timer.C:
			timer.Stop()
			doorTimeout <- true
		}
	}

}

func PollFloors() <-chan int {
	arrivedFloor := make(chan int)
	oldFloor := hw.GetFloorSensorSignal()
	go func() {
		for {
			newFloor := hw.GetFloorSensorSignal()
			if newFloor != oldFloor && newFloor != -1 {
				hw.SetFloorIndicator(newFloor)
				arrivedFloor <- newFloor
			}

			oldFloor = newFloor
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return arrivedFloor
}

func PollButtons() <-chan config.OrderInfo {
	buttonPress := make(chan config.OrderInfo)
	go func() {

		for {
			for floor := 0; floor < config.N_FLOORS; floor++ {
				for button := 0; button < config.N_BUTTONS; button++ {
					if (floor == 0 && button == config.BUTTON_DOWN) || (floor == config.N_FLOORS-1 && button == config.BUTTON_UP) {
						continue
					}

					if hw.GetButtonSignal(button, floor) {
						
							buttonPress <- config.OrderInfo{Button: button, Floor: floor}
							hw.SetButtonLamp(button, floor, true)
						}
						
		
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return buttonPress		
}

