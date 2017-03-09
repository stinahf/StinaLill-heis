package liftAssigner

import (
	//"../eventManager"
	"../config"
	"../queue"
	"fmt"
)

var fittedLiftId string
var minDistance int
var numFittedLifts int
var numLiftsOnNet int
var distances [1][1] int //Husk å endre fra harkoding av antall heiser, lag map

func CalculateBestFit(floor int, button int) {
	numLiftsOnNet = getNumActiveLifts()
	numFittedLifts = 0
	minDistance = 5 
	distance :=  floor - config.InfoPackage[config.IP].CurrentFloor 
	fmt.Println("Jeg gikk inn i CalculateBestFit")

	calculateDistance(floor)

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

	if numFittedLifts < 1 {
		for lift := 0; lift < numLiftsOnNet; lift ++{

			if config.InfoPackage[distances[lift][0]].State == config.Idle {

				if distances[lift][1] < minDistance{
					minDistance = distances[lift][i]
				}
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
		fmt.Println(fittedLiftId)

	}
	
	if numFittedLifts < 1{
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i alle fordi ingen passet kravene")
	}

	for id/*, elevatorInfo*/ := range config.InfoPackage{
		info := queue.OrderInfo{false, id, nil}
		if id != fittedLiftId && fittedLiftId != config.IP{
			fmt.Println(id)
			fmt.Println(fittedLiftId)
			queue.AddSafetyOrder(floor, button, info)
			fmt.Println("La til i safetyOrder")
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
	
func calculateDistance(floor int) {
	var distance int
	for id := range config.InfoPackage {
		distance = floor - config.InfoPackage[id].CurrentFloor
		if distance < 1{
			distance = -1*distance
		}
		distances[id][distance]
	}
}

func getNumActiveLifts() int{
	for id := range config.InfoPackage {
		numLiftsOnNet++
	}
	return numLiftsOnNet
}

