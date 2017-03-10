package liftAssigner

import (
	"../config"
	"../queue"
	"fmt"
	//"time"
)

var fittedLiftId string
var minDistance int
var numLiftsOnNet int
var numFittedLifts int
var liftAssigned bool

func Init() {
	config.Distances = make(map[string]config.DistanceInfo)
}

func BestLift(floor int, button int) {
	numLiftsOnNet = getNumActiveLifts()
	numFittedLifts = 0
	minDistance = 5
	distance := floor - config.InfoPackage[config.IP].CurrentFloor
	fmt.Println("Jeg gikk inn i CalculateBestFit")

	calculateDistance(floor)

	for id := range config.InfoPackage {
		if numFittedLifts < 1 {
			switch button {
			case config.BUTTON_UP:
				if config.InfoPackage[id].MotorDir == config.DIR_UP {
					config.Distances[id] = config.DistanceInfo{id, distance, true}
				}
				if config.InfoPackage[id].MotorDir == config.DIR_UP && distance < 3 {
					queue.AddLocalOrder(floor, button)
					fmt.Println("La til i heis med retning oppover")
					numFittedLifts++
					fittedLiftId = id
				}
			case config.BUTTON_DOWN:
				if config.InfoPackage[id].MotorDir == config.DIR_DOWN {
					config.Distances[id] = config.DistanceInfo{id, distance, true}
				}
				if config.InfoPackage[id].MotorDir == config.DIR_DOWN && distance > -3 {
					queue.AddLocalOrder(floor, button)
					fmt.Println("La til i heis med retning nedover")
					numFittedLifts++
					config.Distances[id] = config.DistanceInfo{id, -1 * distance, true}
					fittedLiftId = id
				}
			}
		}
	}

	for id := range config.Distances {
		if config.Distances[id].GotOrder == true {
			liftAssigned = true
		}
	}

	if numFittedLifts < 1 && !liftAssigned {
		for id := range config.Distances {
			if config.InfoPackage[config.Distances[id].Id].State == config.Idle || config.InfoPackage[config.Distances[id].Id].State == config.DoorOpen && queue.IsQueueEmpty() {

				if config.Distances[id].Distance < minDistance {
					minDistance = config.Distances[id].Distance
					config.Distances[id] = config.DistanceInfo{id, minDistance, true}
					fmt.Println("min distance", config.Distances[id])
				}
			}
		}
	}

	if minDistance == config.Distances[config.IP].Distance {
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i heis i idle og med minste distanse")
		numFittedLifts++
		fmt.Println(numFittedLifts)
		fittedLiftId = config.IP
		fmt.Println(fittedLiftId)
	}

	for id := range config.Distances {
		if config.Distances[id].GotOrder == true {
			liftAssigned = true
		}
	}

	if numFittedLifts < 1 && !liftAssigned {
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i alle fordi ingen passet kravene")

	}

	for id := range config.InfoPackage {
		info := queue.OrderInfo{false, id, nil}
		if id != fittedLiftId && fittedLiftId != config.IP {
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
		config.Distances[id] = config.DistanceInfo{id, distance, false}
		fmt.Println(config.Distances[id])
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
