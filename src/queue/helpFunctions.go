package queue

import (
	"../config"
	"../hw"
	"fmt"
	"time"
)

func (q *queue) setOrder(floor int, button int, status OrderInfo) {
	if q.matrix[floor][button].Active == status.Active {
		return
	}
	fmt.Println("SetOrder: ", status)
	q.matrix[floor][button] = status
	hw.SetButtonLamp(button, floor, true)

	newOrder <- true
}

func (q *queue) startTimer(floor, button int) {
	q.matrix[floor][button].Timer = time.NewTimer(time.Second * 30)
	<-q.matrix[floor][button].Timer.C

	message <- config.Message{OrderComplete: false, Floor: floor, Button: button}
}

func (q *queue) stopTimer(floor, button int) {
	if q.matrix[floor][button].Timer != nil {
		q.matrix[floor][button].Timer.Stop()
	}
}

func (q *queue) isQueueEmpty() bool {
	for floor := 0; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if q.matrix[floor][button].Active {
				return false
			}
		}
	}
	fmt.Println("The queue is empty - sleepy time! :D")
	return true
}

func (q *queue) isOrderAbove(currentFloor int) bool {
	for floor := currentFloor + 1; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if q.matrix[floor][button].Active {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrderBelow(currentFloor int) bool {
	for floor := 0; floor < currentFloor; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if q.matrix[floor][button].Active {
				return true
			}
		}
	}
	return false
}

func (q *queue) shouldStop(dir int, floor int) bool {
	switch dir {
	case config.DIR_UP:
		return q.matrix[floor][config.BUTTON_UP].Active || q.matrix[floor][config.BUTTON_INTERNAL].Active || floor == config.N_FLOORS-1 || !q.isOrderAbove(floor)
	case config.DIR_DOWN:
		return q.matrix[floor][config.BUTTON_DOWN].Active || q.matrix[floor][config.BUTTON_INTERNAL].Active || floor == 0 || !q.isOrderBelow(floor)
	case config.DIR_STOP:
		return q.matrix[floor][config.BUTTON_DOWN].Active || q.matrix[floor][config.BUTTON_INTERNAL].Active || q.matrix[floor][config.BUTTON_UP].Active
	}
	return false
}

func (q *queue) chooseMotorDirection(floor int, dir int) int {
	if q.isQueueEmpty() {
		return config.DIR_STOP
	}
	switch dir {
	case config.DIR_DOWN:
		if q.isOrderBelow(floor) && floor > 0 {
			return config.DIR_DOWN
		} else {
			return config.DIR_UP
		}
	case config.DIR_UP:
		if q.isOrderAbove(floor) && floor < config.N_FLOORS-1 {
			return config.DIR_UP
		} else {
			return config.DIR_DOWN
		}
	case config.DIR_STOP:
		if q.isOrderAbove(floor) {
			return config.DIR_UP
		} else if q.isOrderBelow(floor) {
			return config.DIR_DOWN
		} else {
			return config.DIR_STOP
		}
	}
	return 0
}
