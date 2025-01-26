package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/shashwatrathod/redis-internals/config"
	"github.com/shashwatrathod/redis-internals/core"
)

var concurrent_clients = 0

// Starts a simple TCP server that accepts incoming connections
// and handles commands.
func RunSyncTcpServer() {
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

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		cmd, err := readCommand(conn)

		if err != nil {
			conn.Close()
			concurrent_clients -= 1
			log.Println("Client ", conn.RemoteAddr(), " disconnected. Concurrent clients = ", concurrent_clients)
			if err == io.EOF {
				break
			}
			log.Println("Error", err)
		}

		respond(cmd, conn)
		if err != nil {
			log.Println("Error while sending response: ", err)
		}
	}
}

// reads a single RESP-encoded command from the connection, decodes it,
// and returns a `RedisCmd`.
func readCommand(c io.ReadWriter) (*core.RedisCmd, error) {
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

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

// Decodes the provided RESP-encoded bytes into a
// slice of decoded string tokens.
func decodeArrayString(data []byte) ([]string, error) {
	decodedVal, err := core.Decode(data)

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

func respond(cmd *core.RedisCmd, c io.ReadWriter) {
	err := core.EvalAndRespond(cmd, c)

	if err != nil {
		encodedError := core.EncodeResp(err, false)
		c.Write(encodedError)
	}
}
