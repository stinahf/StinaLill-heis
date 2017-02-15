package queue

import (
	"./config"
    "eventManager"
)

type queue struct {
    matrix [def.N_FLOORS][def.N_BUTTONS]orderInfo
    }

type orderInfo struct {
    active bool
    elev_id int
    timer
}

var newLocalOrder = make(chan bool)
var OrderTimeout = make(chan def.OrderInfo)
var newOrder = make(chan bool)

var local_queue queue
var safety_queue queue

func QueueInit() //hva bør denne gjøre?

func AddLocalOrder(floor int, button int, info orderInfo){
    if local_queue.matrix[floor][button].active == info.active{
        return
    }
    local_queue.matrix[floor][button].active = true
    newOrder <- true
}

func AddSafetyOrder(floor int, button int, info orderInfo){
    if isExternalOrder(button){
        if safety_queue.matrix[floor][button].active == info.active{
            return
        }
        safety_queue.matrix[floor][button].active = true
    }
    return
}

func RemoveOrder(floor int, Message chan<- def.Message){
    for button = 0; b < N_BUTTONS; button++{
        local_queue.matrix[floor][b].active = false
        safety_queue.matrix[floor][b].active = false
        //somethingstoptimeronsafetyqueueorders
    }
    Message <- def.Message(Status: def.OrderComplete, Floor: floor)

}

func RemoveSafetyOrder(floor int, button int, info orderInfo){
    for button = 0; b < N_BUTTONS; button++{
        safety_queue.matrix[floor][b].active = false
        //somethingstoptimeronsafetyqueueorders
    }
}

func isExternalOrder(button int) bool{
    if button == BUTTON_INTERNAL{
        return false
    }
    return true
}
