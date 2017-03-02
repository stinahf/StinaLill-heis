package queue

import (
	"../config"
    "fmt"

)

type orderInfo struct {
    active bool
    elev_id int
    timer bool
}

type queue struct {
    matrix [config.N_FLOORS][config.N_BUTTONS]orderInfo
    }

var newLocalOrder chan bool
var OrderTimeout chan config.OrderInfo
var NewOrder chan bool
//var newOrder 

var local_queue queue
var safety_queue queue

func PrintMatrix() {
    for f := config.N_FLOORS-1; f>=0; f-- {
        print1 := ""
        if local_queue[f][config.BUTTON_UP].active {
            print1 += "↑"
        } else {
            print1 += " "
        }
        if local_queue[f][config.BUTTON_INTERNAL].active {
            print1 += "x"
        } else {
            print1 += " "
        }
        fmt.Println(print1)
        if local_queue[f][config.BUTTON_DOWN].active {
            fmt.Println("↓   %d  ", f+1)
        } else {
            fmt.Println("    %d  ", f+1)
        }
    }
}

func GetFloorFromQueue() int{
    return local_queue[0][0]
}

func GetButtonFromQueue() int{
    return local_queue[0][1]
}


func Init() {
    newLocalOrder = make(chan bool)
    OrderTimeout = make(chan config.OrderInfo)
    NewOrder = make(chan bool)


    //ch.newOrder = make(chan bool)

    fmt.Println("Queue initialized")

}

func AddLocalOrder(floor int, button int){
    if local_queue.matrix[floor][button].active == true {
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

func RemoveOrder(floor int, Message chan<- config.Message){
    for button := 0; button < config.N_BUTTONS; button++{
        local_queue.matrix[floor][button].active = false
        safety_queue.matrix[floor][button].active = false
        //somethingstoptimeronsafetyqueueorders
    }
    Message <- config.Message{Status: config.OrderComplete, Floor: floor}

}

func RemoveSafetyOrder(floor int, info orderInfo){
    for button := 0; button < config.N_BUTTONS; button++{
        safety_queue.matrix[floor][button].active = false
        //somethingstoptimeronsafetyqueueorders
    }
}

func isExternalOrder(button int) bool{
    if button == config.BUTTON_INTERNAL{
        return false
    }
    return true
}

func (q queue) isQueueEmpty() bool{
    for floor := 0; floor < config.N_FLOORS; floor++ {
        for button := 0; button < config.N_BUTTONS; button++ {
            if local_queue.matrix[floor][button].active {
                return false 
            }
        }
    }
    return true
}

func (q queue) isOrderAbove(currentFloor int) bool{
    for floor := currentFloor + 1; floor < config.N_FLOORS; floor++ {
        for button := 0; button < config.N_BUTTONS; button++ {
            if local_queue.matrix[floor][button].active {
                return true
            }
        }
    }
    return false 
}

func (q queue) isOrderBelow(currentFloor int) bool{
    for floor := 0; floor < currentFloor; floor++ {
        for button := 0; button < config.N_BUTTONS; button++ {
            if local_queue.matrix[floor][button].active {
                return true
            }
        }
    }
    return false
}

func ShouldStop(dir int, floor int) bool { 
    switch dir {
    case config.DIR_UP:
        return local_queue.matrix[floor][config.BUTTON_UP].active || local_queue.matrix[floor][config.BUTTON_INTERNAL].active
    case config.DIR_DOWN:
        return local_queue.matrix[floor][config.BUTTON_DOWN].active || local_queue.matrix[floor][config.BUTTON_INTERNAL].active
    case config.DIR_STOP:
        return local_queue.matrix[floor][config.BUTTON_DOWN].active || local_queue.matrix[floor][config.BUTTON_INTERNAL].active || local_queue.matrix[floor][config.BUTTON_UP].active

    }
    return false
}

func ChooseMotorDirection(floor int, dir int) int {
    if local_queue.isQueueEmpty(){
        return config.DIR_STOP
    }
    switch dir{
    case config.DIR_DOWN:
        if local_queue.isOrderBelow(floor){
            return config.DIR_DOWN
        } else {
            return config.DIR_UP
        }
    case config.DIR_UP:
        if local_queue.isOrderAbove(floor){
            return config.DIR_UP
        } else{
            return config.DIR_DOWN
        }
    case config.DIR_STOP:
        if local_queue.isOrderAbove(floor){
            return config.DIR_UP
        } else if local_queue.isOrderBelow(floor){
            return config.DIR_DOWN
        } else{
            return config.DIR_STOP
        }
    }
    return 0
}
