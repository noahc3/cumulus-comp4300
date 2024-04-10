package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"middleware/httpnode"
	"middleware/socketnode"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/akamensky/argparse"
)

// all middleware node types can connect and send messages to cloud servers, regardless of implementation
type Node interface {
	GetIp() string
	GetPort() int
	ConnectToServer() error
	SendToServer(string) (interface{}, error)
	ReadFromServer() (interface{}, error)
}

func main() {
	parser := *argparse.NewParser("middleware", "Create a node between the client and cloud server")
	serverType := parser.String("t", "type", nil)
	serverAddress := parser.String("s", "server-address", nil)
	nodeIp := parser.String("i", "node-ip", nil)
	nodePort := parser.Int("", "node-port", nil)
	connProtocol := parser.String("p", "protocol", nil)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var node Node
	sType := strings.ToLower(*serverType)
	if sType == "websocket" {
		node = httpnode.CreateHttpNode(*serverAddress, *nodeIp, *nodePort)
	} else if sType == "api" {
		fmt.Println("Not yet implemented")
	} else if sType == "socket" {
		node = socketnode.CreateSocketNode(*serverAddress, *nodeIp, *nodePort)
	}

	err = node.ConnectToServer()
	if err != nil {
		fmt.Println("Error:", err)
	}

	cProtocol := strings.ToLower(*connProtocol)
	if cProtocol == "tcp" {
		err = ListenTcp(node)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else if cProtocol == "udp" {
		err = ListenUdp(node)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		err = ListenWs(node)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func ListenTcp(node Node) error {
	tcpAddress := net.TCPAddr{
		IP:   net.ParseIP(node.GetIp()),
		Port: node.GetPort(),
	}

	nodeTcpConn, err := net.ListenTCP("tcp", &tcpAddress)
	if err != nil {
		return err
	}

	for {
		clientConn, err := nodeTcpConn.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		println("Connected")

		var response interface{}
		response, _ = node.ReadFromServer()
		fmt.Println("Received:", response)

		responseString, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
		}

		// send the response back to the client
		_, err = clientConn.Write(responseString)
		if err != nil {
			fmt.Println("Error:", err)
		}

		go func(clientConn net.Conn) {
			for {
				message, err := bufio.NewReader(clientConn).ReadString('\n')
				if err != nil {
					fmt.Println(err)
					break
				}

				fmt.Println("Received:", message)

				response, err := node.SendToServer(message)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}

				fmt.Println("Received:", response)

				responseString, err := json.Marshal(response)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}

				// send the response back to the client
				_, err = clientConn.Write(responseString)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}
			}
		}(clientConn)

	}
}

func ListenUdp(node Node) error {
	udpAddress := net.UDPAddr{
		IP:   net.ParseIP(node.GetIp()),
		Port: node.GetPort(),
	}
	nodeUdpConn, err := net.ListenUDP("udp", &udpAddress)
	if err != nil {
		return err
	}
	defer nodeUdpConn.Close()

	for {
		buffer := make([]byte, 1204)
		_, _, err := nodeUdpConn.ReadFromUDP(buffer[0:])
		if err != nil {
			return err
		}

		// buffer will have a lot of nulls unless it's 1024 character; trim these for readability
		trimmedJson := strings.Trim(string(buffer), "\x00")
		fmt.Println("Received (client):", trimmedJson)

		var v map[string]any
		err = json.Unmarshal([]byte(trimmedJson), &v)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// type assertion: make sure the value held in "value" is a string
		msg, isString := v["value"].(string)
		if !isString {
			return errors.New("value provided by client was not a string")
		}

		response, err := node.SendToServer(msg)
		if err != nil {
			return err
		}

		fmt.Println("Received (server):", response)
		responseString, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
		}

		// argument is a float by default; we just make sure it's a number we can downcast to an int
		portString, isString := v["returnAddr"].(string)
		if !isString {
			return errors.New("return address provided by client was not a number")
		}

		port, _ := strconv.Atoi(portString)

		trimmedMsg := strings.Trim(string(responseString), "\n\x00")
		_, err = nodeUdpConn.WriteToUDP([]byte(trimmedMsg), &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: int(port)})
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func ListenWs(node Node) error {
	ctx := context.WithoutCancel(context.Background())

	listenerAddr := fmt.Sprintf("%v:%v", node.GetIp(), node.GetPort())
	l, err := net.Listen("tcp", listenerAddr)
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	// assert node is of type SocketNode so it implements ServeHTTP
	socketNode, _ := node.(*socketnode.SocketNode)

	server := &http.Server{
		Handler:      socketNode,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	// create another thread that enqueues any server errors into the channel errc
	errc := make(chan error, 1)
	go func() {
		errc <- server.Serve(l)
	}()

	// wait for user-inputted commands…

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	return server.Shutdown(ctx)
}
