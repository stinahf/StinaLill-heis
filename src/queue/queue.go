package queue

import (
	def "config"
    "eventManager"
)

type orderInfo struct {
    active bool
    elev_id int
    timer
}

type queue struct {
    matrix [def.N_FLOORS][def.N_BUTTONS]orderInfo
    }

var newLocalOrder
var OrderTimeout 
//var newOrder 

var local_queue queue
var safety_queue queue

func Init() {
    ch.newLocalOrder = make(chan bool)
    ch.OrderTimeout = make(chan def.OrderInfo)
    //ch.newOrder = make(chan bool)

}

func AddLocalOrder(floor int, button int, info orderInfo){
    if local_queue.matrix[floor][button].active == info.active{
        return
    }
    local_queue.matrix[floor][button].active = true
    NewOrder <- true
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

func isQueueEmpty() bool{
    for floor := 0; floor < def.N_FLOORS; floor++{
        for button := 0, button < def.N_BUTTONS; buttons++{
            if local_queue.matrix[floor][button].active{
                return false 
            }

        }
    }
    return true 
}
func isOrderAbove(currentFloor int) bool{
    for floor := currentFloor + 1; floor < def.N_FLOORS; floor++{
        for button := 0, button < def.N_BUTTONS; buttons++{
            if local_queue.matrix[floor][button].active{
                return true
            }

        }
    }
    return false 
}

func isOrderBelow(currentFloor int) bool{
    for floor := 0; floor < currentFloor; floor++{
        for button := 0, button < def.N_BUTTONS; buttons++{
            if local_queue.matrix[floor][button].active{
                return true
            }

        }
    }
    return false
}

func ShouldStop(dir int, floor int) bool { 
    switch dir {
    case def.DIR_UP:
        return local_queue.matrix[floor][BUTTON_UP].active || local_queue.matrix[floor][BUTTON_INTERNAL].active
    case def.DIR_DOWN:
        return local_queue.matrix[floor][BUTTON_DOWN].active || local_queue.matrix[floor][BUTTON_INTERNAL].active
    case def.DIR_STOP:
        local_queue.matrix[floor][BUTTON_DOWN].active || local_queue.matrix[floor][BUTTON_INTERNAL].active || local_queue.matrix[floor][BUTTON_UP].active

    }
}

func ChooseMotorDirection(floor int, dir int) int {
    if local_queue.IsQueueEmpty{
        return def.DIR_STOP
    }
    switch dir{
    case def.DIR_DOWN:
        if local_queue.IsOrderBelow(floor){
            return def.DIR_DOWN
        }
        else {
            return def.DIR_UP
        }
    case def.DIR_UP:
        if local_queue.IsOrderAbove(floor){
            return def.DIR_UP
        }
        else{
            return def.DIR_DOWN
        }
    case def.DIR_STOP:{
        if local_queue.IsOrderAbove(floor){
            return def.DIR_UP
        }
        else if local_queue.IsOrderBelow(floor){
            return def.DIR_DOWN
        }
        else{
            return def.DIR_STOP
        }
    }

}