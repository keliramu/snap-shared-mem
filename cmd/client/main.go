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
	log.Println("~~~Client START~~~")
	log.Println("~~~SNAP_COMMON:", os.Getenv("SNAP_COMMON"))
	log.Println("~~~SNAP_DATA:", os.Getenv("SNAP_DATA"))
	u, _ := user.Current()
	log.Println("~~~current user:", u)
	log.Println("~~~socket file:", internal.SockAddrSh)

	// write file in shared memory
	data := []byte("Hello from Client!\n")
	err1 := os.WriteFile(internal.JustFile, data, 0644)
	log.Println("~~~WriteFile err:", err1)

	// connect to server on unix socket in shared memory
	var conn net.Conn
	var err error

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)

		conn, err = net.Dial(internal.Protocol, internal.SockAddrSh)
		if err != nil {
			log.Fatal(err)
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
