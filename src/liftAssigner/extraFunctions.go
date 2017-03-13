package liftAssigner

import (
	"../config"
	"../queue"
)

func initGotOrder() {
	for id := range config.InfoPackage {
		config.GotOrder[id] = config.GotOrderInfo{id, false}
	}
}

func initBestLift(floor int) (int, int){
	minDistance = 5
	liftAssigned = false
	smallestId = config.InfoPackage[config.IP].Id
	distance := floor - config.InfoPackage[config.IP].CurrentFloor

	return smallestId, distance
}


func assignIfIdle(floor int, button int) {
	for id := range config.Distances {
		if config.InfoPackage[config.Distances[id].Id].State == config.Idle || config.InfoPackage[config.Distances[id].Id].State == config.DoorOpen && queue.IsQueueEmpty() {
			if config.Distances[id].Distance < minDistance {
				minDistance = config.Distances[id].Distance
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
				}
			}
		if smallestId == config.Distances[config.IP].Id {
			liftAssigned = true
			queue.AddLocalOrder(floor, button)
			bestLiftId = config.IP
		}
	}
}

func assignIfMoving(floor int, button int, distance int) {
	for id := range config.InfoPackage {
		switch button {
		case config.BUTTON_UP:
			if config.InfoPackage[id].MotorDir == config.DIR_UP && distance < 0 {
				if config.Distances[id].Distance < minDistance {
					minDistance = config.Distances[id].Distance
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
					}
				}
					if smallestId == config.Distances[config.IP].Id {
						liftAssigned = true
						queue.AddLocalOrder(floor, button)
						bestLiftId = config.IP
					}
			}
		case config.BUTTON_DOWN:
			if config.InfoPackage[id].MotorDir == config.DIR_DOWN && distance > 0{
				if config.Distances[id].Distance < minDistance {
					minDistance = config.Distances[id].Distance
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
					}
				}
					if smallestId == config.Distances[config.IP].Id {
						liftAssigned = true
						queue.AddLocalOrder(floor, button)
						bestLiftId = config.IP
					}
			}
		}
	}
}

func assignIfNoBestFit(floor int, button int) {
	for id := range config.InfoPackage {
		if config.InfoPackage[id].Id > smallestId {
				smallestId = config.Distances[id].Id
			}
		}
	if smallestId == config.Distances[config.IP].Id {
		liftAssigned = true
		queue.AddLocalOrder(floor, button)
		bestLiftId = config.IP
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


func updateLiftAssigned() {
		for id := range config.GotOrder {
		if config.GotOrder[id].GotOrder == true {
			liftAssigned = true
		}
	}
}


