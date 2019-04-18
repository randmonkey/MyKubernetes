package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type RES struct {
	Host string `json: "host"`
	Ip   string `json: "ip"`
	Time string `json: "time"`
}

type MRES struct {
	Res     RES
	Version string `json: "version`
}

func main() {
	http.HandleFunc("/status", statusRespons)
	http.HandleFunc("/more", moreResponse)
	http.ListenAndServe(":80", nil)
}

func moreResponse(w http.ResponseWriter, req *http.Request) {
	var mRes MRES
	mRes.Res = getLocal()
	mRes.Version = "v0.2.0"
	r, _ := json.Marshal(mRes)
	w.Write(r)

}

func statusRespons(w http.ResponseWriter, req *http.Request) {
	r, _ := json.Marshal(getLocal())
	w.Write(r)
}

func getLocal() (localstatus RES) {
	var ipv4 string
	nicname := getInterface()
	addrs, err := net.InterfaceByName(nicname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ips, _ := addrs.Addrs()
	for _, ip := range ips {
		if len(ip.String()) < 19 {
			ipv4 = ip.String()
		}
	}

	localstatus.Host, _ = os.Hostname()
	localstatus.Ip = ipv4
	localstatus.Time = time.Now().String()
	return
}

func getInterface() string {
	route, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		fmt.Println("error of get route file")
	}
	lines := strings.Split(string(route), "\n")
	for _, line := range lines {
		words := strings.Split(line, "\t")
		var r []string
		for _, word := range words {
			r = append(r, word)
		}
		if len(r) > 0 {
			if r[1] == "00000000" && r[7] == "00000000" {
				return r[0]
			}
		}
	}
	return ""
}
