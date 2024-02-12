package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"snapshared/internal"
	"time"
)

func main() {
	//sockAddrInSnap := os.Getenv("SNAP_COMMON")
	log.Println("~~~Client START~~~")
	log.Println("~~~SNAP_COMMON:", os.Getenv("SNAP_COMMON"))
	log.Println("~~~SNAP_DATA:", os.Getenv("SNAP_DATA"))
	u, _ := user.Current()
	log.Println("~~~current user:", u)

	sockAddrInSnap := internal.SockAddrSh //os.Getenv("SNAP_COMMON") + sockAddr

	log.Println("~~~socket file:", sockAddrInSnap)

	var conn net.Conn
	var err error

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)

		conn, err = net.Dial(internal.Protocol, sockAddrInSnap)
		if err != nil {
			log.Println("~~~err1:", err)
			conn, err = net.Dial(internal.Protocol, internal.SockClAddrSh)
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = conn.Write([]byte("hello world"))
		if err != nil {
			log.Fatal(err)
		}

		err = conn.(*net.UnixConn).CloseWrite()
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
