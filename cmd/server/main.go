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
// https://snapcraft.io/blog/private-shared-memory-support-for-snaps
// https://forum.snapcraft.io/t/the-shared-memory-interface/28382
// https://snapcraft.io/docs/shared-memory-interface

func main() {
	log.Println("~~~Server START~~~")
	log.Println("~~~SNAP_COMMON:", os.Getenv("SNAP_COMMON"))
	log.Println("~~~SNAP_DATA:", os.Getenv("SNAP_DATA"))
	u, _ := user.Current()
	log.Println("~~~current user:", u)
	log.Println("~~~socket file:", internal.SockAddrSh)

	cleanup := func() {
		if _, err := os.Stat(internal.SockAddrSh); err == nil {
			if err := os.RemoveAll(internal.SockAddrSh); err != nil {
				log.Fatal(err)
			}
		}
	}

	cleanup()

	// write file in shared memory
	data := []byte("Hello from Server!\n")
	err := os.WriteFile(internal.JustFile, data, 0644)
	log.Println("~~~WriteFile err:", err)

	// server on unix socket in shared memory
	listener, err := net.Listen(internal.Protocol, internal.SockAddrSh)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	quit := make(chan os.Signal, 1)
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
