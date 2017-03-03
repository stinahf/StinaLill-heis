package hw

import (
	"../config"
	"errors"
	"fmt"
	"log"
)

//const MOTOR_SPEED 2800 //TODO - Move to config file

var lamp_channel_matrix = [config.N_FLOORS][config.N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [config.N_FLOORS][config.N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Init() error {
	if !ioInit() {
		return errors.New("Hw: ioInit() failed!")
	}

	for f := 0; f < config.N_FLOORS; f++ {
		if f != 0 {
			SetButtonLamp(config.BUTTON_DOWN, f, false)
		}
		if f != config.N_FLOORS-1 {
			SetButtonLamp(config.BUTTON_UP, f, false)
		}
		SetButtonLamp(config.BUTTON_INTERNAL, f, false)
	}

	SetStopLamp(false)
	SetDoorOpenLamp(false)

	//Move to init floor (1st floor)
	SetMotorDirection(config.DIR_DOWN)
	floor := GetFloorSensorSignal()
	for floor == -1 {
		floor = GetFloorSensorSignal()
	}
	SetMotorDirection(config.DIR_STOP)
	SetFloorIndicator(floor)

	fmt.Println("Hw initialized")
	return nil

}

func SetMotorDirection(dirn int) {
	if dirn == 0 {
		ioWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		ioClearBit(MOTORDIR)
		ioWriteAnalog(MOTOR, 2800)
	} else if dirn < 0 {
		ioSetBit(MOTORDIR)
		ioWriteAnalog(MOTOR, 2800)
	}
}

func SetButtonLamp(button int, floor int, value bool) {
	if floor < 0 || floor >= config.N_FLOORS {
		log.Printf("Error: The floor is out of range", floor)
		return
	}
	if button == config.BUTTON_UP && floor == config.N_FLOORS-1 {
		log.Println("You are already at the top")
		return
	}
	if button == config.BUTTON_DOWN && floor == 0 {
		log.Println("You are already at the bottom")
		return
	}
	if button != config.BUTTON_UP && button != config.BUTTON_DOWN && button != config.BUTTON_INTERNAL {
		log.Printf("Invalid button %d\n", button)
		return
	}

	if value {
		ioSetBit(lamp_channel_matrix[floor][button])
	} else {
		ioClearBit(lamp_channel_matrix[floor][button])
	}
}

func SetFloorIndicator(floor int) {
	if floor < 0 || floor >= config.N_FLOORS {
		log.Printf("The floor %d is out of range! \n", floor)
		return
	}

	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
}

func SetDoorOpenLamp(value bool) {
	if value {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

func SetStopLamp(value bool) {
	if value {
		ioSetBit(LIGHT_STOP)
	} else {
		ioClearBit(LIGHT_STOP)
	}
}

func GetButtonSignal(button int, floor int) bool {
	if floor < 0 || floor >= config.N_FLOORS {
		log.Printf("The floor %d is out of range \n", floor)
		return false
	}
	if button < 0 || button >= config.N_BUTTONS {
		log.Printf("Button %d is out of range \n", floor)
		return false
	}
	if button == config.BUTTON_UP && floor == config.N_FLOORS-1 {
		log.Println("You are already on the top")
		return false
	}
	if button == config.BUTTON_DOWN && floor == 0 {
		log.Println("You are already at the buttom")
		return false
	}
	if ioReadBit(button_channel_matrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func GetFloorSensorSignal() int {
	if ioReadBit(SENSOR_FLOOR1) {
		return 0
	} else if ioReadBit(SENSOR_FLOOR2) {
		return 1
	} else if ioReadBit(SENSOR_FLOOR3) {
		return 2
	} else if ioReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func GetStopSignal() bool {
	return ioReadBit(STOP)
}

func GetObstructionSignal() bool {
	return ioReadBit(OBSTRUCTION)
}