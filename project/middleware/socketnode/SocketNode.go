package socketnode

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type SocketNode struct {
	ServerAddress string
	NodeIp        string
	NodePort      int
	ctx           context.Context
	clientConn    *websocket.Conn
	serverConn    net.Conn
}

func CreateSocketNode(serverAddress string, nodeIp string, nodePort int) *SocketNode {
	sn := &SocketNode{
		ServerAddress: serverAddress,
		NodeIp:        nodeIp,
		NodePort:      nodePort,
	}
	return sn
}

func (sn *SocketNode) GetIp() string { return sn.NodeIp }

func (sn *SocketNode) GetPort() int { return sn.NodePort }

func (sn *SocketNode) ConnectToServer() error {
	// connect to the supplied address
	var err error
	sn.serverConn, err = net.Dial("tcp", sn.ServerAddress)
	if err != nil {
		return err
	}

	return nil
}

func (sn *SocketNode) SendToServer(msg string) (interface{}, error) {
	fmt.Println("Sending " + msg)
	_, err := sn.serverConn.Write([]byte(msg))
	if err != nil {
		return nil, err
	}

	fmt.Println("Sent")

	// wait for the server to respond
	response := make([]byte, 1024)
	_, err = sn.serverConn.Read(response)
	if err != nil {
		return nil, err
	}

	fmt.Println("Received from server:", string(response))

	return string(response), nil
}

func (sn *SocketNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept websocket connections from clients
	var err error
	sn.clientConn, err = websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println("Error:", err)
		sn.clientConn.CloseNow()
	}
	defer sn.clientConn.CloseNow()

	sn.ctx = context.WithoutCancel(r.Context())

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for {
		var v string
		err = wsjson.Read(sn.ctx, sn.clientConn, &v)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				fmt.Println("Client closed connection")
				break
			}

			fmt.Println("Error reading:", err)
			return
		} else {
			fmt.Println("Received:", v)
		}

		response, err := sn.SendToServer(v)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println("Sending to client:", response)

		err = wsjson.Write(sn.ctx, sn.clientConn, response)
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}

func (sn *SocketNode) ReadFromServer() (interface{}, error) {
	// to implement
	return nil, nil
}
