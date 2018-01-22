package main

import (
	"fmt"
	//"net/http"
	//"regexp"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"path"
)

//var data string = `node_bonding_active{master="bond10"} 1`

func getNicStatus() {
	//var ifnames []string
	links, err := netlink.LinkList()
	if err != nil {
		fmt.Errorf("Get link list err :", err.Error())
	}
	for _, link := range links {
		if link.Type() == "bond" {
			ifname := link.Attrs().Name
			filename := path.Join("/sys/class/net/", ifname, "speed")
			buf, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Errorf("Error open bond speed file", err.Error())
			}
			fmt.Println(ifname, string(buf))
		}
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
