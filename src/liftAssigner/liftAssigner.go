package liftAssigner

import (
	"../config"
	"../hw"
	"../queue"
	"fmt"
	//"time"
)

var fittedLiftId int
var minDistance int
var numLiftsOnNet int
var numFittedLifts int
var liftAssigned bool
var smallestId int

func Init() {
	config.Distances = make(map[int]config.DistanceInfo)
	config.GotOrder = make(map[int]config.GotOrderInfo)
}

func BestLift(floor int, button int) {
	hw.SetButtonLamp(button, floor, true)
	numLiftsOnNet = getNumActiveLifts()
	numFittedLifts = 0
	minDistance = 5
	liftAssigned = false
	smallestId = config.InfoPackage[config.IP].Id
	distance := floor - config.InfoPackage[config.IP].CurrentFloor
	initGotOrder()
	fmt.Println(config.GotOrder[config.IP])
	fmt.Println("Jeg gikk inn i CalculateBestFit")
	fmt.Println("init liftAssigned: ", liftAssigned)
	calculateDistance(floor)

	for id := range config.InfoPackage {
		if config.InfoPackage[id].State == config.DoorOpen && distance == 0 {
			queue.AddLocalOrder(floor, button)
			config.GotOrder[id] = config.GotOrderInfo{id, true}
		}
	}

	for id := range config.GotOrder {
		if config.GotOrder[id].GotOrder == true {
			liftAssigned = true
			fmt.Println("InitAssigned", liftAssigned)
		}
	}

	if !liftAssigned {
		for id := range config.Distances {
			if config.InfoPackage[config.Distances[id].Id].State == config.Idle || config.InfoPackage[config.Distances[id].Id].State == config.DoorOpen && queue.IsQueueEmpty() {
				fmt.Println("Jeg kom meg nesten til minDistance, men mangler litt")
				if config.Distances[id].Distance < minDistance {
					minDistance = config.Distances[id].Distance
					fmt.Println("min distance", config.Distances[id])
				}
			}
		}
		for id := range config.Distances {
			if config.Distances[id].Distance == minDistance {
				config.GotOrder[id] = config.GotOrderInfo{id, true}
			}
		}

		if minDistance == config.Distances[config.IP].Distance {
			for id := range config.Distances {
				if config.Distances[id].Id < smallestId && config.GotOrder[id].GotOrder {
					smallestId = config.Distances[id].Id
					fmt.Println("smallest ID: ", smallestId)
				}
			}
			if smallestId == config.Distances[config.IP].Id {
				fmt.Println("Jeg har minst ID", smallestId)
				liftAssigned = true
				queue.AddLocalOrder(floor, button)
				fittedLiftId = config.IP
			}
		}

	}

	for id := range config.GotOrder {
		if config.GotOrder[id].GotOrder == true {
			liftAssigned = true
			fmt.Println("LiftAssigned", liftAssigned, "after checking if idle")
		}
	}
	if !liftAssigned {
		for id := range config.InfoPackage {
			switch button {
			case config.BUTTON_UP:
				if config.InfoPackage[id].MotorDir == config.DIR_UP && distance > 0 {
					fmt.Println("Gikk inn første if med retning opp")
					config.GotOrder[id] = config.GotOrderInfo{id, true}
				}
				if config.InfoPackage[id].MotorDir == config.DIR_UP && distance < 0 {
					queue.AddLocalOrder(floor, button)
					fmt.Println("La til i heis med retning oppover")
					fittedLiftId = id
				}
			case config.BUTTON_DOWN:
				if config.InfoPackage[id].MotorDir == config.DIR_DOWN {
					fmt.Println("Gikk inn i retning ned")
					config.GotOrder[id] = config.GotOrderInfo{id, true}
				}
				if config.InfoPackage[id].MotorDir == config.DIR_DOWN {
					queue.AddLocalOrder(floor, button)
					fmt.Println("La til i heis med retning nedover")
					fittedLiftId = id
				}
			}
		}
	}

	/*if numFittedLifts < 1 && !liftAssigned {
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i alle fordi ingen passet kravene")

	}*/

	for id := range config.InfoPackage {
		info := queue.OrderInfo{true, id, nil}
		if id != fittedLiftId && fittedLiftId != config.IP {
			queue.AddSafetyOrder(floor, button, info)
			fmt.Println("La til i safetyOrder")
		}
	}
}

func HandleExternalOrderStatus(msgInfo config.Message) {
	if msgInfo.OrderComplete == true {
		queue.RemoveSafetyOrder(msgInfo.Floor)
	} else {
		queue.AddLocalOrder(msgInfo.Floor, msgInfo.Button)
	}
}

func calculateDistance(floor int) {
	var distance int
	for id := range config.InfoPackage {
		distance = floor - config.InfoPackage[id].CurrentFloor
		if distance < 1 {
			distance = -1 * distance
			fmt.Println("Distansen kalkulert: ", distance)
		}
		config.Distances[id] = config.DistanceInfo{id, distance}
		fmt.Println(config.Distances[id])
	}
}

func initGotOrder() {
	for id := range config.InfoPackage {
		config.GotOrder[id] = config.GotOrderInfo{id, false}
	}
}

func getNumActiveLifts() int {
	for id := range config.InfoPackage {
		if id != config.InfoPackage[id].Id {
			numLiftsOnNet++
		}
	}
	return numLiftsOnNet
}
