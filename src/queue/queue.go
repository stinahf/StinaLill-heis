package queue

import (
	"../config"
	"../hw"
	"fmt"
	"time"
)

type OrderInfo struct {
	Active  bool
	Elev_id string      `json:"-"`
	Timer   *time.Timer `json:"-"`
}

type queue struct {
	matrix [config.N_FLOORS][config.N_BUTTONS]OrderInfo
}

var newOrder chan bool
var newLocalOrder chan bool
var OrderTimeout chan OrderInfo
var message chan config.Message

var local_queue queue
var safety_queue queue

func Init(newOrderTemp chan bool, messageTemp chan config.Message) {
	message = messageTemp
	newOrder = newOrderTemp
	newLocalOrder = make(chan bool)
	OrderTimeout = make(chan OrderInfo)

	fmt.Println("Queue initialized")

}

func AddLocalOrder(floor int, button int) {
	local_queue.setOrder(floor, button, OrderInfo{true, " ", nil})
}

func AddSafetyOrder(floor int, button int, info OrderInfo) {
	safety_queue.setOrder(floor, button, info)
	go safety_queue.startTimer(floor, button)
}

func RemoveOrder(floor int) {
	for button := 0; button < config.N_BUTTONS; button++ {
		local_queue.matrix[floor][button].Active = false
		safety_queue.matrix[floor][button].Active = false
		hw.SetButtonLamp(button, floor, false)
	}
	message <- config.Message{OrderComplete: true, Floor: floor, Button: 0}
}

func RemoveSafetyOrder(floor int) {
	for button := 0; button < config.N_BUTTONS; button++ {
		safety_queue.matrix[floor][button].Active = false
		safety_queue.stopTimer(floor, button)
	}
}

func isExternalOrder(button int) bool {
	if button == config.BUTTON_INTERNAL {
		return false
	}
	return true
}

func IsQueueEmpty() bool {
	return local_queue.isQueueEmpty()
}

func IsOrderAbove(currentFloor int) bool {
	return local_queue.isOrderAbove(currentFloor)
}

func IsOrderBelow(currentFloor int) bool {
	return local_queue.isOrderBelow(currentFloor)
}

func ShouldStop(dir int, floor int) bool {
	return local_queue.shouldStop(dir, floor)
}

func ChooseDirection(floor int, dir int) int {
	fmt.Println(local_queue.chooseMotorDirection(floor, dir))
	return local_queue.chooseMotorDirection(floor, dir)
}
