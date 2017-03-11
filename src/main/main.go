package main

import (
	"../Network"
	"../config"
	"../eventManager"
	"../hw"
	"../liftAssigner"
	"../queue"
	"fmt"
)

func main() {

	ch := eventManager.Channels{
		NewOrder:       make(chan bool),
		ReachedFloor:   make(chan int),
		MotorDir:       make(chan int),
		DoorLamp:       make(chan bool),
		DoorTimerReset: make(chan bool),
		DoorTimeout:    make(chan bool),
	}

	channels := Network.ReceiveChannels{
		ReceiveMessage:       make(chan config.Message),
		ReceiveInfo:          make(chan config.ElevatorMsg),
		ReceiveExternalOrder: make(chan config.OrderInfo),
	}

	Network.Message = make(chan config.Message)
	NewExternalOrder := make(chan config.OrderInfo)

	hw.Init()
	eventManager.Init()
	queue.Init(ch.NewOrder, Network.Message)
	Network.Init()
	liftAssigner.Init()

	sendInfo := Network.SendInfoPacket()

	go Network.Transmitter(16569, Network.Message)
	go Network.Receiver(16569, channels.ReceiveMessage)

	go Network.Transmitter(16571, NewExternalOrder)
	go Network.Receiver(16571, channels.ReceiveExternalOrder)

	go Network.Transmitter(16570, sendInfo)
	go Network.Receiver(16570, channels.ReceiveInfo)

	go Network.NetworkHandler(channels)

	go eventManager.EventManager(ch)
	go eventManager.OpenDoor(ch.DoorTimeout, ch.DoorTimerReset)

	manageEvents(ch, NewExternalOrder, channels)

}

func manageEvents(ch eventManager.Channels, New chan config.OrderInfo, channels Network.ReceiveChannels) {
	buttonPress := make(chan config.OrderInfo)
	go eventManager.PollButtons(buttonPress)
	floorHIT := make(chan int)
	go eventManager.PollFloors(floorHIT)
	for {
		select {
		case button := <-buttonPress:
			switch button.Button {
			case config.BUTTON_UP, config.BUTTON_DOWN:
				New <- button
			case config.BUTTON_INTERNAL:
				queue.AddLocalOrder(button.Floor, button.Button)
			}
		case floor := <-floorHIT:
			ch.ReachedFloor <- floor
		case dir := <-ch.MotorDir:
			hw.SetMotorDirection(dir)
		case value := <-ch.DoorLamp:
			hw.SetDoorOpenLamp(value)
		case messageInfo := <-channels.ReceiveMessage:
			fmt.Println("OrderComplete: ", messageInfo)
			liftAssigner.HandleExternalOrderStatus(messageInfo)
		}
	}
}
