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
	Id           int
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
	Button        int
}

const (
	Idle     = 0
	Moving   = 1
	DoorOpen = 2
)

type DistanceInfo struct {
	Id       int
	Distance int
}

type GotOrderInfo struct {
	Id       int
	GotOrder bool
}

var InfoPackage map[int]ElevatorMsg
var Distances map[int]DistanceInfo
var GotOrder map[int]GotOrderInfo

var IP int
