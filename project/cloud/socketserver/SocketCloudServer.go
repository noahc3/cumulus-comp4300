package socketserver

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type SocketCloudServer struct {
}

func CreateSocketServer() *SocketCloudServer {
	scs := &SocketCloudServer{}
	return scs
}

func (scs SocketCloudServer) HandleTcp(clientConn net.Conn) {
	// connection should remain open as long as the user has more data to send
	for {
		message, err := bufio.NewReader(clientConn).ReadString('\n')
		if err != nil && strings.Contains(err.Error(), "EOF") {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println("Received:", message)

		// create a json response and send it back to the client, trimming nulls and newlines
		trimmedMessage := strings.Trim(message, "\n")
		trimmedMessage = strings.Trim(trimmedMessage, "\x00")
		messageJson := fmt.Sprintf("{\"message\": \"%v\"}", trimmedMessage)
		_, err = clientConn.Write([]byte(messageJson))
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
