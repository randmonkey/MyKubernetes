package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("start tcp server")
	l, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("error of listen port")
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on 0.0.0.0:9999 ...")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error of accept content")
			os.Exit(1)
		}
		defer conn.Close()
		fmt.Printf("Receive %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}
