package liftAssigner

import (
	//"../eventManager"
	"../config"
	"../queue"
	"fmt"
)

var fittedLiftId string
var numLiftsIdle int



func CalculateBestFit(floor int, button int) {
	numFittedLifts := 0
	minDistance := 5 
	distance :=  floor - config.InfoPackage[config.IP].CurrentFloor 
	fmt.Println("Jeg gikk inn i CalculateBestFit")

	for id := range config.InfoPackage {
		switch button {
		case config.BUTTON_UP:
			if config.InfoPackage[id].MotorDir == config.DIR_UP && distance > 0  {
				queue.AddLocalOrder(floor, button)
				fmt.Println("La til i heis med retning oppover")
				numFittedLifts++
				fittedLiftId = id
			}
		case config.BUTTON_DOWN:
			if config.InfoPackage[id].MotorDir == config.DIR_DOWN && distance < 0 {
				queue.AddLocalOrder(floor, button)
				fmt.Println("La til i heis med retning nedover")
				numFittedLifts++
				fittedLiftId = id
			}
		}
	}

	for id := range config.InfoPackage{
		if numFittedLifts < 1 && config.InfoPackage[id].State == config.Idle {

			if distance > 0{
				if (distance) < minDistance{
					minDistance = distance
				}
			}

			fmt.Println("IsOrderBelow: ", distance)
			if distance < 0{
				distance = -1*distance
				if distance < minDistance{
					fmt.Println("Setting minDistance")
					minDistance = distance
				}
			}
		}


	//if 	   (queue.IsOrderAbove(config.InfoPackage[id].CurrentFloor) && minDistance == distanceOrderAbove) 
	//	|| (queue.IsOrderBelow(config.InfoPackage[id].CurrentFloor) && minDistance == distanceOrderBelow) {
	if 	minDistance == distance {
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i heis i idle og med minste distanse")
		numFittedLifts ++
		fittedLiftId = config.IP

	}
	
	if numFittedLifts < 1{
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i alle fordi ingen passet kravene")
	}

	for id/*, elevatorInfo*/ := range config.InfoPackage{
		info := queue.OrderInfo{true, id, nil}
		if id != fittedLiftId{
			fmt.Println(id)
			fmt.Println(fittedLiftId)
			queue.AddSafetyOrder(floor, button, info)
			fmt.Println("La til i safetyOrder")
			}
		}
	}
}


func HandleExternalOrderStatus(msgInfo config.Message) {
	if msgInfo.OrderComplete == true {
		queue.RemoveSafetyOrder(msgInfo.Floor)
	}else {
		queue.AddLocalOrder(msgInfo.Floor, msgInfo.Button)
	}
}
	
