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

func Init() {
	config.Distances = make(map[string]config.DistanceInfo)
}

func CalculateBestFit(floor int, button int) {
	numLiftsOnNet = getNumActiveLifts()
	numFittedLifts = 0
	minDistance = 5
	distance := floor - config.InfoPackage[config.IP].CurrentFloor
	fmt.Println("Jeg gikk inn i CalculateBestFit")

	calculateDistance(floor)

	for id := range config.InfoPackage {
		switch button {
		case config.BUTTON_UP:
			if config.InfoPackage[id].MotorDir == config.DIR_UP && distance > 0 {
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
		for id := range config.Distances {

			if config.InfoPackage[config.Distances[id].Id].State == config.Idle {

				if config.Distances[id].Distance < minDistance {
					minDistance = config.Distances[id].Distance
					fmt.Println(minDistance)
				}
			}
		}
	}

	if minDistance == distance {
		queue.AddLocalOrder(floor, button)
		fmt.Println("La til i heis i idle og med minste distanse")
		numFittedLifts++
		fittedLiftId = config.IP
		fmt.Println(fittedLiftId)

	}

	if numFittedLifts < 1 {
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
		}
		config.Distances[id] = config.DistanceInfo{id, distance}
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

//FANT BUG! TAR ALDRI HENSYN TIL OM DØRA ER ÅPEN!
