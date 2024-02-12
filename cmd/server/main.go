package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"snapshared/internal"
	"strings"
)

// https://eli.thegreenplace.net/2019/unix-domain-sockets-in-go/
func main() {
	log.Println("~~~Server START~~~")
	log.Println("~~~SNAP_COMMON:", os.Getenv("SNAP_COMMON"))
	log.Println("~~~SNAP_DATA:", os.Getenv("SNAP_DATA"))
	u, _ := user.Current()
	log.Println("~~~current user:", u)

	sockAddrInSnap := internal.SockAddrSh //os.Getenv("SNAP_COMMON") + sockAddr

	log.Println("~~~socket file:", sockAddrInSnap)

	cleanup := func() {
		if _, err := os.Stat(sockAddrInSnap); err == nil {
			if err := os.RemoveAll(sockAddrInSnap); err != nil {
				log.Fatal(err)
			}
		}
	}

	cleanup()

	listener, err := net.Listen(internal.Protocol, sockAddrInSnap)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("ctrl-c pressed..")
		close(quit)
		cleanup()
		os.Exit(0)
	}()

	fmt.Println("server launched...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(">>> accepted")
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connected: %s\n", conn.RemoteAddr().Network())

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, conn)
	if err != nil {
		log.Println(err)
		return
	}

	s := strings.ToUpper(buf.String())

	buf.Reset()
	buf.WriteString(s)

	_, err = io.Copy(conn, buf)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("<<< ", s)
}
