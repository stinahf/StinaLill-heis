package network

import (
	"../queue"
	"../liftAssigner"
	"../config"
	"../eventManager"
	"net"
	"strconv"
	"strings"
	"time"
)

func NetworkHandler(ch ReceiveChannels) {
	for {
		select {
		case receiveElevMsg := <-ch.ReceiveInfo:
			receiveInfoPacket(receiveElevMsg)
		case receiveExternal := <-ch.ReceiveExternalOrder:
			liftAssigner.BestLift(receiveExternal.Floor, receiveExternal.Button)
			queue.AddSafetyOrder(receiveExternal.Floor, receiveExternal.Button, queue.OrderInfo{true, nil})
		}
	}
}

func receiveExternalOrder(receiveExternal config.OrderInfo) {
	config.ExternalOrderInfo = config.OrderInfo{receiveExternal.Button, receiveExternal.Floor}
}


func receiveInfoPacket(receivePacket config.ElevatorMsg) {
	config.InfoPackage[receivePacket.Id] = config.ElevatorMsg{receivePacket.Id, receivePacket.CurrentFloor, receivePacket.MotorDir, receivePacket.State}
	infoPackageTimer[receivePacket.Id] = time.Now().Unix()
	for id := range infoPackageTimer {
		if time.Now().Unix() - infoPackageTimer[id] >= 2 {
			mutex.Lock()
			delete(config.InfoPackage, id)
			delete(config.GotOrder, id)
			delete(config.Distances, id)
			mutex.Unlock()
		}
	}
}

func SendInfoPacket() <-chan config.ElevatorMsg {
	sendInfo := make(chan config.ElevatorMsg)
	ip := getIP()
	config.IP = splitIP(ip)
	go func() {
		for {
			Floor, Dir, State := eventManager.GetFloorDirState()
			sendPacket := config.ElevatorMsg{Id: config.IP, CurrentFloor: Floor, MotorDir: Dir, State: State}
			sendInfo <- sendPacket
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return sendInfo
}

func getIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "" 
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue 
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue 
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "" 
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue 
			}
			return ip.String()
		}
	}
	return "" 
}

func splitIP(IP string) int {
	ip := strings.Split(IP, ".")[3]
	id, _ := strconv.ParseInt(ip, 10, 0)
	return int(id)
}