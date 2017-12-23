package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/u35s/rudp"
)

func read(conn *rudp.RudpUnConn) {
	go func() {
		for {
			conn.Tick <- 1
			time.Sleep(1e6)
		}
	}()
	for {
		data := make([]byte, rudp.MAX_PACKAGE)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Printf("read err %s\n", err)
			break
		}
		for i := range data[:n] {
			v := int(data[i])
			if v == 255 {
				fmt.Printf("receive ")
				fmt.Printf("%d", v)
				fmt.Printf(" from <%v>\n", conn.RemoteAddr())
			}
		}
	}
}

func main() {
	log.SetOutput(os.Stdout)
	addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 9981}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	listener := rudp.NewListener(conn)
	defer func() { fmt.Println("defer close", listener.Close()) }()
	go func() {
		for {
			rconn, err := listener.AcceptRudp()
			if err != nil {
				fmt.Printf("accept err %v\n", err)
				break
			}
			go read(rconn)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT)
	select {
	case <-signalChan:
	}
}
