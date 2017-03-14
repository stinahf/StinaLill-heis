package queue

import (
	"../config"
	"../hw"
	"fmt"
	"time"
)

type OrderInfo struct {
	Active  bool
	Timer   *time.Timer `json:"-"`
}

type queue struct {
	matrix [config.N_FLOORS][config.N_BUTTONS]OrderInfo
}

var newOrder chan bool
var message chan config.Message

var local_queue queue
var safety_queue queue

func Init(newOrderTemp chan bool, messageTemp chan config.Message) {
	message = messageTemp
	newOrder = newOrderTemp

	fmt.Println("Queue initialized")
	loadFromHardware(local_queue)
	filterInternalQueue(local_queue)

}

func filterInternalQueue(q queue) {
	for floor:=0; floor<config.N_FLOORS; floor++ {
		for button:= 0; button<config.BUTTON_INTERNAL; button++ {
			q.matrix[floor][button].Active = false
		}
	}
}


func AddLocalOrder(floor int, button int) {
	local_queue.setOrder(floor, button, OrderInfo{true, nil})
	newOrder <- true
}

func AddSafetyOrder(floor int, button int, info OrderInfo) {
	safety_queue.setOrder(floor, button, info)
	go safety_queue.startTimer(floor, button)
}

func RemoveOrder(floor int) {
	for button := 0; button < config.N_BUTTONS; button++ {
		hw.SetButtonLamp(button, floor, false)
		local_queue.matrix[floor][button].Active = false
		message <- config.Message{OrderComplete: true, Floor: floor, Button: button}
	}
}

func RemoveSafetyOrder(floor int) {
	for button := 0; button < config.BUTTON_INTERNAL; button++ {
		hw.SetButtonLamp(button, floor, false)
		safety_queue.stopTimer(floor, button)
		safety_queue.matrix[floor][button].Active = false
	}
}

func IsQueueEmpty() bool {
	return local_queue.isQueueEmpty()
}

func ShouldStop(dir int, floor int) bool {
	return local_queue.shouldStop(dir, floor)
}

func ChooseDirection(floor int, dir int) int {
	return local_queue.chooseMotorDirection(floor, dir)
}

func ShouldAddOrder(floor int, button int) bool {
	for button := 0; button < config.BUTTON_INTERNAL; button++ {
	if safety_queue.matrix[floor][button].Active {
		return false
	}
}
	return true
}
