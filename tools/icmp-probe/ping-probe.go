package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	ping "github.com/sparrc/go-ping"
)

func main() {
	ipList := os.Getenv("IPLIST")
	//ipList := []string{"183.136.237.66", "127.0.0.1", "172.217.160.100"}
	fmt.Println("starting a ping probe...")
	fmt.Println("targets is:", ipList)
	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		ips := strings.Split(ipList, " ")
		for _, target := range ips {
			go doPing(target)
		}
	}
}

func doPing(target string) {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		fmt.Printf("error of ping %s, %v", target, err)
	}
	pinger.Count = 5
	pinger.Interval, _ = time.ParseDuration("0.2s")
	pinger.Timeout, _ = time.ParseDuration("2s")
	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	if stats.PacketLoss != 0 {
		fmt.Println(time.Now().Format("2006-01-02T15:04:05Z07:00"), stats.Addr, stats.PacketsSent, stats.PacketsRecv, stats.MaxRtt)
	}
}
