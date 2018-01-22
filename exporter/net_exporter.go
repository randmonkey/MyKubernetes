package main

import (
	"fmt"
	//"net/http"

	"github.com/vishvananda/netlink"
)

//var data string = `node_bonding_active{master="bond10"} 1`

func getNicStatus() {
	links, err := netlink.LinkList()
	if err != nil {
		fmt.Errorf("Get link list err :", err.Error())
	}
	for _, link := range links {
		fmt.Println(link.Attrs().Name)
	}
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, data)
// }

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe("0.0.0.0:9001", nil)
	getNicStatus()
}
