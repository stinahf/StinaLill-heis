package queue

import (
	"../config"
	"../hw"
	"fmt"
	"time"
)

type queue struct {
	matrix [config.N_FLOORS][config.N_BUTTONS]orderInfo
}


var newOrder chan bool
var newLocalOrder chan bool
var OrderTimeout chan config.OrderInfo


//var newOrder

var local_queue queue
var safety_queue queue


func Init(newOrderTemp chan bool) {
	newOrder = newOrderTemp
	newLocalOrder = make(chan bool)
	OrderTimeout = make(chan config.OrderInfo)

	//ch.newOrder = make(chan bool)

	fmt.Println("Queue initialized")

}

func (q *queue) setOrder(floor int, button int, status orderInfo) {
	if q.matrix[floor][button].active == status.active {
		return
	}
	q.matrix[floor][button] = status
	hw.SetButtonLamp(button, floor, true)

	newOrder <- true
}

func AddLocalOrder(floor int, button int, id int) {
	local_queue.setOrder(floor, button, orderInfo{true/*, id*/, nil})
}

func AddSafetyOrder(floor int, button int, info orderInfo) {
	if isExternalOrder(button) {
		if safety_queue.matrix[floor][button].active == info.active {
			return
		} else {
			safety_queue.setOrder(floor, button, orderStatus{true, nil})
			go safety_queue.startTimer(floor, button)
			
		}
	}
	return
}

func (q *queue) startTimer(floor, button int) {
	q.matrix[floor][button].timer = time.NewTimer(time.Second * 30)
	<-q.matrix[floor][button].timer.C
	OrderTimeout <- OrderInfo{Floor: floor, Button: button}
}

func (q *queue) stopTimer(floor, button int) {
	if q.matrix[floor][button].timer != nil {
		q.matrix[floor][button].timer.Stop()
	}
}

func RemoveOrder(floor int) { //Husk, ta inn: , Message chan<- config.Message
	for button := 0; button < config.N_BUTTONS; button++ {
		local_queue.matrix[floor][button].active = false 
		safety_queue.matrix[floor][button].active = false
		hw.SetButtonLamp(button, floor, false)
	}
	//Message <- config.Message{Status: config.OrderComplete, Floor: floor}
}

func RemoveSafetyOrder(floor int, info orderInfo) {
	for button := 0; button < config.N_BUTTONS; button++ {
		safety_queue.matrix[floor][button].active = false
		safety_queue.stopTimer(floor, button)
	}
}

func isExternalOrder(button int) bool {
	if button == config.BUTTON_INTERNAL {
		return false
	}
	return true
}

func (q *queue) IsQueueEmpty() bool {
	for floor := 0; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if q.matrix[floor][button].active {
				return false
				fmt.Println("The queue is not empty, get your ass moving, someone is waiting!")
			}
		}
	}
	fmt.Println("The queue is empty - sleepy time! :D")
	return true
}

func IsQueueEmpty() bool{
	return local_queue.IsQueueEmpty()
}

func (q *queue) isOrderAbove(currentFloor int) bool {
	for floor := currentFloor + 1; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if q.matrix[floor][button].active {
				return true
			}
		}
	}
	return false
}

func IsOrderAbove(currentFloor int) bool {
	return local_queue.isOrderAbove(currentFloor)
}

func (q *queue) isOrderBelow(currentFloor int) bool {
	for floor := 0; floor < currentFloor; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if q.matrix[floor][button].active {
				return true
			}
		}
	}
	return false
}

func IsOrderBelow(currentFloor int) bool {
	return local_queue.isOrderBelow(currentFloor)
}

func (q *queue) shouldStop(dir int, floor int) bool {
	switch dir {
	case config.DIR_UP:
		return q.matrix[floor][config.BUTTON_UP].active || q.matrix[floor][config.BUTTON_INTERNAL].active || floor == config.N_FLOORS-1 || !q.isOrderAbove(floor)
	case config.DIR_DOWN:
		return q.matrix[floor][config.BUTTON_DOWN].active || q.matrix[floor][config.BUTTON_INTERNAL].active || floor == 0 || !q.isOrderBelow(floor)
	case config.DIR_STOP:
		return q.matrix[floor][config.BUTTON_DOWN].active || q.matrix[floor][config.BUTTON_INTERNAL].active || q.matrix[floor][config.BUTTON_UP].active
	}
	return false
}

func ActuallyShouldStop(dir int, floor int) bool {
	return local_queue.shouldStop(dir, floor)
}

func (q *queue) ChooseMotorDirection(floor int, dir int) int {
	if q.IsQueueEmpty() {
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

func ActuallyChooseDirection(floor int, dir int) int {
	fmt.Println(local_queue.ChooseMotorDirection(floor, dir))
	return local_queue.ChooseMotorDirection(floor, dir)
}
