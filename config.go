package config

const (
	BUTTON_UP = 0
	BUTTON_DOWN = 1
	BUTTON_INTERNAL = 2
)

const N_FLOORS = 4
const N_BUTTONS = 3

const MOTOR_SPEED = 2800

const (
	DIR_UP = 1
	DIR_STOP = 0
	DIR_DOWN = -1
)

type ElevatorInfo struct {
	CurrentFloor int
	MotorDir int
	State //If we opt to have states
	//Queue something something
}

type NewOrderInfo struct {
	Button int
	Floor int
}

const (
	Idle = 0
	Moving = 1
	OpenDoor = 2
)

var state int