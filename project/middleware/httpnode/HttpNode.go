package httpnode

import (
	"context"
	"fmt"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type HttpNode struct {
	ServerAddress string
	NodeIp        string
	NodePort      int
	ctx           context.Context
	conn          *websocket.Conn
}

func CreateHttpNode(serverAddress string, nodeIp string, nodePort int) *HttpNode {
	hn := &HttpNode{
		ServerAddress: serverAddress,
		NodeIp:        nodeIp,
		NodePort:      nodePort,
	}
	return hn
}

func (hn *HttpNode) GetIp() string { return hn.NodeIp }

func (hn *HttpNode) GetPort() int { return hn.NodePort }

func (hn *HttpNode) ConnectToServer() error {
	hn.ctx = context.WithoutCancel(context.Background())

	protocolString := fmt.Sprintf("ws://%v", hn.ServerAddress)
	// dial into the cloud server websocket
	var err error
	hn.conn, _, err = websocket.Dial(hn.ctx, protocolString, nil)
	if err != nil {
		return err
	}

	// start the process
	fmt.Println("Starting process…")
	startProcessMap := make(map[string]string)
	startProcessMap["command"] = "start"
	startProcessMap["directory"] = "./"
	startProcessMap["target"] = "../../mockprocess/echo.out"
	return wsjson.Write(hn.ctx, hn.conn, startProcessMap)

}

func (hn *HttpNode) SendToServer(msg string) (interface{}, error) {
	// forward the message to our server
	msgTrimmed := strings.Trim(msg, "\n\x00")
	msgJson := make(map[string]string)
	msgJson["command"] = "input"
	msgJson["value"] = msgTrimmed

	err := wsjson.Write(hn.ctx, hn.conn, msgJson)
	if err != nil {
		fmt.Println("Error writing")
		return nil, err
	}

	fmt.Println("Sending to server:", msgJson)

	var v interface{}
	err = wsjson.Read(hn.ctx, hn.conn, &v)
	if err != nil {
		return nil, err
	}

	fmt.Println("Received from server:", v)

	return v, err
}

func (hn *HttpNode) ReadFromServer() (interface{}, error) {
	var v interface{}
	err := wsjson.Read(hn.ctx, hn.conn, &v)
	if err != nil {
		return nil, err
	}

	fmt.Println("Received from server:", v)

	return v, err
}
