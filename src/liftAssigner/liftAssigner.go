package liftAssigner

import (
	"../eventManager"
	"../config"
	"../queue"
)

var numFittedLifts int


func liftAssigner(floor, button) {
	for id, elevatorInfo := range infoPackage {
		switch button {
		case BUTTON_UP:
			if dir == DIR_UP && IsOrderAbove(infoPackage.floor) && (floor - infoPackage[id].floor) < 3 {
				queue.AddLocalQueue(floor, button, id)
				numFittedLifts++
			}
		case BUTTON_DOWN:
			if dir == DIR_DOWN && IsOrderBelow(infoPackage.floor) && (infoPackage[id].floor - floor) < 3{
				queue.AddLocalQueue(floor, button, id)
				numFittedLifts++
			}
		}
	}

	for id, elevatorInfo := infoPackage{
		if numFittedLifts < 1{
		
			if infoPackage[id].state == idle{
				queue.AddLocalQueue(floor, button, id)
				numFittedLifts++;
				
			}
		}
	}
	
	if numFittedLifts < 1{
		queue.AddLocalQueue(floor, button, infoPackage[id].id)
	}

		}
	}
}
	
