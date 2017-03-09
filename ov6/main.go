package main

import (
	"./network"
	"fmt"
	"time"
)


func main(){
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	controllCh := make(chan int)
	broadcastCh := make(chan int)
	master := false
	go network.Network(controllCh, broadcastCh, &master)
	var mynumber int
	t0 := time.Now()

	
	mynumber = 0
	go func(){
		for {
		number := network.GetNumber()
		d := time.Since(t0)
		if (d.Seconds() > 5){
			if (number == 0){
				master = true
				println("master is true")
				
			} 
			
		}
		if (number > mynumber){
			mynumber = number
		}
		println(mynumber)
		mynumber++
		network.SendMsg(broadcastCh, mynumber)
		time.Sleep(1000*time.Millisecond)

		}
		}()

	time.Sleep(100*time.Second)
	


	
	

}