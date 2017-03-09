package config

const (
	BUTTON_UP       = 0
	BUTTON_DOWN     = 1
	BUTTON_INTERNAL = 2
)

const N_FLOORS = 4
const N_BUTTONS = 3

const MOTOR_SPEED = 2800

const (
	DIR_UP   = 1
	DIR_STOP = 0
	DIR_DOWN = -1
)

type ElevatorMsg struct {
	Id           string
	CurrentFloor int
	MotorDir     int
	State        int
}

type OrderInfo struct {
	Button int
	Floor  int
}

var ExternalOrderInfo OrderInfo

type Message struct {
	OrderComplete bool
	Floor         int
	Button 		  int
}

const (
	Idle        = 0
	Moving      = 1
	DoorClosing = 2
)

var InfoPackage map[string]ElevatorMsg

var IP string
