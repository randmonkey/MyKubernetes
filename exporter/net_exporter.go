package main

import (
	"fmt"
	//"net/http"

	"github.com/vishvananda/netlink"
)

//var data string = `node_bonding_active{master="bond10"} 1`

func getNicStatus() {
	link, err := netlink.LinkList()
	if err != nil {
		fmt.Errorf("Get link list err :", err.Error())
	}
	ifname := link.Attrs().Name
	fmt.Println(ifname)
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, data)
// }

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe("0.0.0.0:9001", nil)
	getNicStatus()
}
