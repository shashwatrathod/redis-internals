package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"syscall"

	"github.com/shashwatrathod/redis-internals/config"
	"github.com/shashwatrathod/redis-internals/core/commandhandler"
	"github.com/shashwatrathod/redis-internals/core/eval"
	redisio "github.com/shashwatrathod/redis-internals/core/io"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
)

// GREAT video on FDs https://www.youtube.com/watch?v=-gP58pozNuM

const (
	max_concurrent_clients = 20000
)

func RunAsyncTcpServer() error {
	log.Println("Initializing the server on ", config.Host, ":", config.Port)

	// First initialize a socket.
	// O_NONBLOCK creates the socket in a non-blocking mode.
	// SOCK_STREAM sets the type of socket to STREAM.
	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)

	if err != nil {
		return err
	}
	defer syscall.Close(serverFd)

	// Enable nonblocking behavior on the socket server.
	if err = syscall.SetNonblock(serverFd, true); err != nil {
		return err
	}

	// Bind the socket to the given port.
	ipv4 := net.ParseIP(config.Host).To4()
	err = syscall.Bind(serverFd, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]},
	})
	if err != nil {
		return err
	}

	// Start listening on the socket server for new connections.
	//  max_concurrent_clients specifies the max number of clients that can be in the queue.
	if err = syscall.Listen(serverFd, max_concurrent_clients); err != nil {
		return err
	}

	// ONLY FOR LINUX
	// Create a new Epoll through system call.
	// Epoll can be thought of as an "Observable" in the "Observer" pattern
	// that monitors the information and passes that information to the observers.
	epollFd, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}

	// This is the event that is going to be "observed".
	// EPOLLIN gets fired when a file descriptor (in this case serverFd) is ready to be read.
	var socketEvent = &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFd),
	}

	// Adds a new "observer" for the Event. The observer gets attached to Observable (EPOLL)'s FD.
	if err = syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, serverFd, socketEvent); err != nil {
		return err
	}

	log.Println("Sucessfully started the server.")
	log.Printf("Listening on %s:%d...\n", config.Host, config.Port)

	var s store.Store = store.GetStore()

	concurrent_clients := 0

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_concurrent_clients)

	for {
		// Wait for new events to be captured.
		nevents, e := syscall.EpollWait(epollFd, events, -1)
		if e != nil {
			return nil
		}

		for i := 0; i < nevents; i++ {
			var event syscall.EpollEvent = events[i]

			// If there is an event on serverFd,
			// that means a new client wants to connect.
			if event.Fd == int32(serverFd) {
				// Accept the new connection
				conn_fd, conn_address, e := syscall.Accept(serverFd)
				if e != nil {
					log.Println("An error occurred while accepting connection from a client: ", err)
					continue
				}

				concurrent_clients++

				if addr, ok := conn_address.(*syscall.SockaddrInet4); ok {
					ip := net.IPv4(addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3])
					log.Printf("Successfully accepted a connection from %s:%d. Concurrent Clients = %d\n", ip.String(), addr.Port, concurrent_clients)
				}

				if e := syscall.SetNonblock(serverFd, true); e != nil {
					log.Println("Error while configuring the nonblocking server: ", e)
					continue
				}

				// Event where the Client's Fd is ready to be read.
				// This basically means we have new data/information incoming from the client.
				var clientEvent *syscall.EpollEvent = &syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(conn_fd),
				}

				// Add a new "Observer" to listen for events on the Client's FD.
				if e := syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, conn_fd, clientEvent); e != nil {
					log.Println("Error occured while estabilishing listner on Client", err)
					concurrent_clients--
					syscall.Close(conn_fd)
				}
			} else {
				// This means we have a new event on the Client's FD.
				comm := &redisio.FDComm{
					Fd: int(event.Fd),
				}

				command, err := readCommand(comm)
				if err != nil {
					syscall.Close(int(event.Fd))
					concurrent_clients--
					continue
				}

				respond(command, s, comm)
			}
		}
	}

	return nil
}

// reads a single RESP-encoded command from the connection, decodes it,
// and returns a `RedisCmd`.
func readCommand(c io.ReadWriter) (*eval.RedisCmd, error) {
	var buffer []byte = make([]byte, 512)

	size, err := c.Read(buffer)

	if err != nil {
		return nil, err
	}

	if config.LogRequest {
		log.Println("Raw input: ", fmt.Sprintf("%q", buffer[:size]))
	}
	tokens, err := decodeArrayString(buffer[:size])

	if err != nil {
		return nil, err
	}

	return &eval.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

// Decodes the provided RESP-encoded bytes into a
// slice of decoded string tokens.
func decodeArrayString(data []byte) ([]string, error) {
	decodedVal, err := resp.Decode(data)

	if err != nil {
		return []string{}, err
	}

	decodedValues := decodedVal.([]interface{})

	decodedArray := make([]string, len(decodedValues))
	for i, v := range decodedValues {
		decodedArray[i] = v.(string)
	}

	return decodedArray, nil
}

func respond(cmd *eval.RedisCmd, s store.Store, c io.ReadWriter) {
	err := commandhandler.EvalAndRespond(cmd, s, c)

	if err != nil {
		encodedError := resp.Encode(err, false)
		c.Write(encodedError)
	}
}
