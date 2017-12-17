package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func getServerInfo() (ifname string, err error) {
	route, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		err = fmt.Errorf("Error get route, %v", err)
		return
	}
	for _, r := range route {
		if r.Dst == nil {
			link, e := netlink.LinkByIndex(r.LinkIndex)
			if e != nil {
				err = fmt.Errorf("Error get link,%v", e)
				return
			}
			ifname = link.Attrs().Name
			break
		}
	}
	return ifname, err
}

func main() {
	a, _ := getServerInfo()
	fmt.Println(a)
}
