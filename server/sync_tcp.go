package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/shashwatrathod/redis-internals/main/config"
)

var concurrent_clients = 0

// Starts a simple TCP server that accepts incoming connections
// and handles commands.
func RunMultiThreadedSyncTcpServer() {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	log.Println("Starting Simple TCP Server at ", address)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Println("Error starting the TCP server: ", err)
		panic(err)
	}

	log.Println("TCP Server listening for connections at ", address)


	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting incoming connection", err)
			conn.Close()
			panic(err)
		}
		
		concurrent_clients += 1
		log.Println("Client connected with address:", conn.RemoteAddr(), "; concurrent clients = ", concurrent_clients)

		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
	for {
		cmd, err := readCommand(conn)
		
		if (err != nil) {
			conn.Close()
			concurrent_clients -= 1
			log.Println("Client ",conn.RemoteAddr(), " disconnected. Concurrent clients = ", concurrent_clients)
			if (err == io.EOF) {
				break
			}
			log.Println("Error", err)
		}

		log.Println("Command : ", cmd)
		err = echo(cmd, conn)
		if (err != nil) {
			log.Println("Error while sending response: ", err)
		}
	}
}

func echo(cmd string, conn net.Conn) error {
	_, err := conn.Write([]byte(cmd))
	if (err != nil) {
		return err
	}

	return nil
}

func readCommand(conn net.Conn) (string, error) {
	var buffer []byte = make([]byte, 512)

	size, err := conn.Read(buffer)
	
	if (err != nil) {
		return "", err
	}

	return string(buffer[:size]), nil
}