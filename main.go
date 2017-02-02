package main

import (
	"./bcast"
	"fmt"
	"time"
)

type HelloMsg struct {
	Message string
	Iter int
}

func main() {
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	go func() {
		helloMsg := HelloMsg{"Hello from us", 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}