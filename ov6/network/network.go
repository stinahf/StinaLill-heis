package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const SPAMTIME = 1000 //milliseconds


func Network(controllCh chan int, BroadcastCh chan int, master *bool){
	
	
	//port:= "20013"
	//ip:= "255.255.255.255"
	//service :=  fmt.Sprintf("%d:%d", ip, port)
	//service := "129.241.187.255:34767"
	service := "127.255.255.255:12345" //localhost. fungerer ikke. hvordan sette opp ny terminal da?
	addr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		fmt.Printf("Net.ResolveUDPAddr failed!\n")
		return 
	}	

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Printf("Net.DialUDP failed!\n")
		return 
	}	
	recChan = make(chan int)
	

	//broadcastChan := make(chan int)
	//go Broadcast(conn, broadcastChan)

	defer conn.Close()

	localAddr := conn.LocalAddr().String()

	connRec, err := net.ListenUDP("udp", addr)
	if err != nil {
			fmt.Printf("Net.ListenUDP failed!\n")
			return 
		}	

	//hvorfor kan ikke denne deklareres "globalt", slik at receive ikke trenger å ta inn recChan?
	go Receive(connRec, localAddr)
	go Broadcast(conn, BroadcastCh, master)
	time.Sleep(100*time.Millisecond)
	for{


			select {

				
					
				case <-time.After(100*time.Millisecond):

			}

			time.Sleep(100*time.Millisecond)

	}




}

var recChan chan int


func SendMsg(msgChan chan int, msg int) {
	msgChan <- msg
}










func Broadcast(conn net.Conn, broadcastChan chan int, master *bool) {
	// skal sende meldingen vår med et intervall tilsvarende SPAMTIME
	var msg int

	//var delay time.Time 
	for {
		select{
			case msg = <- broadcastChan:
				
		}

		//if time.Since(delay) > SPAMTIME*time.Millisecond { // her kan vi også sjekke om meldingen er valid...
			//delay = time.Now()
				if *master{
					jsonMsg, _ := json.Marshal(msg)
					conn.Write(jsonMsg)	
					}

		//}

	}
}

func GetNumber() int{
	select {
	case curNumber := <- recChan:
		return curNumber
	case <- time.After(100*time.Millisecond):
		return 0
	}

}

func Receive(connRec *net.UDPConn, localAddr string){

	var msg int
	var buf [1024]byte
	for {
		//fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...
		n, _, _ := connRec.ReadFrom(buf[0:])

		//n, receivedAddr, _ := connRec.ReadFrom(buf[0:])
		json.Unmarshal(buf[0:n], &msg)
		//receivedAddr.String() = " " //fjerne denne for å forhindre at meldinger mottas på samme maskin
		
		select {
				case recChan <- msg:
					
	
				case <-time.After(100*time.Millisecond):
					
			}
/*
		if (receivedAddr.String() != localAddr){
			select {
				case recChan <- msg:
					fmt.Println(recChan)
			}
		}
*/
	}
}