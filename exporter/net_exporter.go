package main

import (
	"fmt"
	//"net/http"
	//"regexp"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"path"
	"strings"
)

// var data string = `node_bonding_active{master="bond10"} 1`

func getBondList() (ifname []string) {
	links, err := netlink.LinkList()
	if err != nil {
		fmt.Errorf("Get link list err:", err.Error())
	}
	for _, link := range links {
		if link.Type() == "bond" {
			ifname = append(ifname, link.Attrs().Name)
		}
	}
	return
}

func getBondSpeed(bond string) (speed string) {
	filename := path.Join("/sys/class/net/", bond, "speed")
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("Error open bond speed file", err.Error())
	}
	speed = strings.Replace(string(buf), "\n", "", -1)
	return
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, data)
// }

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe("0.0.0.0:9001", nil)
	bondList := getBondList()
	//fmt.Println(len(bondList))

	for _, bond := range bondList {
		speed := getBondSpeed(bond)
		fmt.Println(bond, speed)
	}
}
