package main

import (
	"../network"
	"../config"
	"../eventManager"
	"../hw"
	"../liftAssigner"
	"../queue"
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

	channels := network.ReceiveChannels{
		ReceiveMessage:       make(chan config.Message),
		ReceiveInfo:          make(chan config.ElevatorMsg),
		ReceiveExternalOrder: make(chan config.OrderInfo),
	}

	network.Message = make(chan config.Message)
	NewExternalOrder := make(chan config.OrderInfo)

	hw.Init()
	eventManager.Init()
	queue.Init(ch.NewOrder, network.Message)
	network.Init()
	liftAssigner.Init()

	sendInfo := network.SendInfoPacket()

	go network.Transmitter(16569, network.Message)
	go network.Receiver(16569, channels.ReceiveMessage)

	go network.Transmitter(16571, NewExternalOrder)
	go network.Receiver(16571, channels.ReceiveExternalOrder)

	go network.Transmitter(16570, sendInfo)
	go network.Receiver(16570, channels.ReceiveInfo)

	go network.NetworkHandler(channels)

	go eventManager.EventManager(ch)
	go eventManager.OpenDoor(ch.DoorTimeout, ch.DoorTimerReset)

	manageEvents(ch, NewExternalOrder, channels)

}

func manageEvents(ch eventManager.Channels, New chan config.OrderInfo, channels network.ReceiveChannels) {
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
			liftAssigner.HandleExternalOrderStatus(messageInfo)
		}
	}
}
