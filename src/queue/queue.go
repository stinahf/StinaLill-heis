package queue

import (
	"../config"
	"fmt"
)

type queue struct {
	matrix [config.N_FLOORS][config.N_BUTTONS]orderInfo
}

type orderInfo struct {
	active bool
	//elev_id int
	//  timer   bool
}

var newLocalOrder chan bool
var OrderTimeout chan config.OrderInfo
var NewOrder chan bool

//var newOrder

var local_queue queue
var safety_queue queue

func PrintMatrix() {
	for f := config.N_FLOORS - 1; f >= 0; f-- {
		print1 := ""
		if local_queue.matrix[f][config.BUTTON_UP].active {
			print1 += "↑"
		} else {
			print1 += " "
		}
		if local_queue.matrix[f][config.BUTTON_INTERNAL].active {
			print1 += "x"
		} else {
			print1 += " "
		}
		fmt.Println(print1)
		if local_queue.matrix[f][config.BUTTON_DOWN].active {
			fmt.Println("↓   %d  ", f+1)
		} else {
			fmt.Println("    %d  ", f+1)
		}
	}
}

func GetFloorFromQueue() orderInfo {
	return local_queue.matrix[0][0]
}

func GetButtonFromQueue() orderInfo {
	return local_queue.matrix[0][1]
}

func Init() {
	newLocalOrder = make(chan bool)
	OrderTimeout = make(chan config.OrderInfo)
	NewOrder = make(chan bool)

	//ch.newOrder = make(chan bool)

	fmt.Println("Queue initialized")

}

func (q *queue) SetOrder(floor int, button int, status orderInfo) {
	fmt.Println("HEIHEIEHIEEI")
	if q.matrix[floor][button].active == status.active {
		return
	}
	fmt.Println("Jaja")
	q.matrix[floor][button] = status
	fmt.Println("Oki, so far so good")

	//NewOrder <- true
}

func AddLocalOrder(floor int, button int/*, id int*/) {
	local_queue.SetOrder(floor, button, orderInfo{true/*, id*/})
}

func AddSafetyOrder(floor int, button int, info orderInfo) {
	if isExternalOrder(button) {
		if safety_queue.matrix[floor][button].active == info.active {
			return
		}
		safety_queue.matrix[floor][button].active = true
	}
	return
}

func RemoveOrder(floor int, Message chan<- config.Message) {
	for button := 0; button < config.N_BUTTONS; button++ {
		local_queue.matrix[floor][button].active = false
		safety_queue.matrix[floor][button].active = false
		//somethingstoptimeronsafetyqueueorders
	}
	Message <- config.Message{Status: config.OrderComplete, Floor: floor}

}

func RemoveSafetyOrder(floor int, info orderInfo) {
	for button := 0; button < config.N_BUTTONS; button++ {
		safety_queue.matrix[floor][button].active = false
		//somethingstoptimeronsafetyqueueorders
	}
}

func isExternalOrder(button int) bool {
	if button == config.BUTTON_INTERNAL {
		return false
	}
	return true
}

func (q *queue) IsQueueEmpty() bool {
	fmt.Println("I'm inside MOHAHA")
	for floor := 0; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if /*local_queue*/ q.matrix[floor][button].active {
				return false
				fmt.Println("The queue is not empty, get your ass moving, someone is waiting!")
			}
		}
	}
	fmt.Println("The queue is empty - sleepy time! :D")
	return true
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
		return q.matrix[floor][config.BUTTON_UP].active || q.matrix[floor][config.BUTTON_INTERNAL].active
	case config.DIR_DOWN:
		return q.matrix[floor][config.BUTTON_DOWN].active || q.matrix[floor][config.BUTTON_INTERNAL].active
	case config.DIR_STOP:
		return q.matrix[floor][config.BUTTON_DOWN].active || q.matrix[floor][config.BUTTON_INTERNAL].active || q.matrix[floor][config.BUTTON_UP].active
	}
	return false
}

func ActuallyShouldStop(dir int, floor int) bool {
	return local_queue.shouldStop(dir, floor)
}

func (q *queue) ChooseMotorDirection(floor int, dir int) int {
	fmt.Println("Lalalallalalalal Lill er bestest!")
	if q.IsQueueEmpty() {
		return config.DIR_STOP
		fmt.Println("Dir stop")
	}
	switch dir {
	case config.DIR_DOWN:
		if q.isOrderBelow(floor) && floor > 0 {
			return config.DIR_DOWN
			fmt.Println("order is below and dir down")
		} else {
			return config.DIR_UP
			fmt.Println("order is above and dir up")
		}
	case config.DIR_UP:
		if q.isOrderAbove(floor) && floor < config.N_FLOORS-1 {
			return config.DIR_UP
			fmt.Println("order is above and dir up")
		} else {
			return config.DIR_DOWN
			fmt.Println("order is below and dir down")
		}
	case config.DIR_STOP:
		fmt.Println("Stina er aller bestest!!!")
		if q.isOrderAbove(floor) {
			return config.DIR_UP
			fmt.Println("dir up")
		} else if q.isOrderBelow(floor) {
			return config.DIR_DOWN
			fmt.Println("dir down")
		} else {
			return config.DIR_STOP
			fmt.Println("dir stop")
		}
	}
	return 0
}

func ActuallyChooseDirection(floor int, dir int) int {
	fmt.Println(local_queue.ChooseMotorDirection(floor, dir))
	return local_queue.ChooseMotorDirection(floor, dir)
}
