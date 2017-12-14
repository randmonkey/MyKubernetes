package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func getServerInfo() {
	link := netlink.NewLinkAttrs()
	link.Name = "enp0s10"
	iflist, err := netlink.AddrList(nil, 4)
	if err != nil {
		fmt.Println("fuck")
	}
	fmt.Println(iflist)
}

func main() {
	getServerInfo()
}
