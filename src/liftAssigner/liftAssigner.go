package liftAssigner

import (
	"../config"
	"../hw"
	"../queue"
)

var bestLiftId int
var minDistance int
var liftAssigned bool
var smallestId int
var distance int

func Init() {
	config.Distances = make(map[int]config.DistanceInfo)
	config.GotOrder = make(map[int]config.GotOrderInfo)
}

func BestLift(floor int, button int) {
	hw.SetButtonLamp(button, floor, true)
	smallestId, distance = initBestLift(floor)
	initGotOrder()
	calculateDistance(floor)

	if queue.ShouldAddOrder(floor, button) {
		for id := range config.InfoPackage {
			if distance == 0 {
				queue.AddLocalOrder(floor, button)
				config.GotOrder[id] = config.GotOrderInfo{id, true}
				bestLiftId = id
			}
		}

		updateLiftAssigned()


		if !liftAssigned {
			assignIfIdle(floor, button)
		}

		updateLiftAssigned()

		if !liftAssigned {
			assignIfMoving(floor, button, distance)
		}

		updateLiftAssigned()

		if !liftAssigned {
			assignIfNoBestFit(floor, button)

		}
	}

}
