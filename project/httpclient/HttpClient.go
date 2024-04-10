package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/akamensky/argparse"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var outputMutex sync.Mutex

func main() {
	parser := argparse.NewParser("httpclient", "Creates a client that can connect to a cloud server over websocket")
	serverAddress := parser.String("s", "server-address", nil)
	testFile := parser.String("t", "test", nil)

	parser.Parse(os.Args)

	var err error
	if *testFile != "" {
		err = runTests(*serverAddress, *testFile)
	} else {
		err = run(*serverAddress)
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Program ended successfully")
}

func run(serverAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	protocolString := fmt.Sprintf("ws://%v/", serverAddress)
	conn, _, err := websocket.Dial(ctx, protocolString, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer conn.CloseNow()

	go readConnOutput(ctx, conn)

	err = startProcess(ctx, conn)
	if err != nil {
		return err
	}

	// wait for user-inputted commands…
	for {
		fmt.Println("I am scanning")
		outputMutex.Lock()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		message, _ := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		outputMutex.Unlock()

		messageJson := make(map[string]string)
		messageJson["command"] = "input"
		messageJson["value"] = message

		var messageString []byte
		messageString, err = json.Marshal(messageJson)
		if err != nil {
			return err
		}

		err = conn.Write(ctx, websocket.MessageText, messageString)
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
	}

	return conn.Close(websocket.StatusNormalClosure, "")
}

func runTests(serverAddress string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	ctx := context.WithoutCancel(context.Background())

	protocolString := fmt.Sprintf("ws://%v/", serverAddress)
	conn, _, err := websocket.Dial(ctx, protocolString, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer conn.CloseNow()

	// read the two first messages
	err = startProcess(ctx, conn)
	if err != nil {
		return err
	}
	readConnOutput(ctx, conn)

	for i := 1; i < 1000; i++ {

		message := "test"
		messageJson := make(map[string]string)
		messageJson["command"] = "input"
		messageJson["value"] = message

		fmt.Println("Writing…")
		err = wsjson.Write(ctx, conn, messageJson)
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
		startTime := time.Now()

		var response []byte
		_, response, err = conn.Read(ctx)
		if err != nil {
			fmt.Println("Error reading response:", err)
		}
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		fmt.Println("Received: ", string(response))

		file.WriteString(fmt.Sprintf("%v\n", float64(latency.Nanoseconds())/float64(10e6)))
	}

	return conn.Close(websocket.StatusNormalClosure, "")
}

func startProcess(ctx context.Context, conn *websocket.Conn) error {
	// start the process
	startProcessMap := make(map[string]string)
	startProcessMap["command"] = "start"
	startProcessMap["directory"] = "./"
	startProcessMap["target"] = "../../mockprocess/echo.out"
	return wsjson.Write(ctx, conn, startProcessMap)
}

func readConnOutput(ctx context.Context, conn *websocket.Conn) {
	for i := 0; i < 1; i++ {
		var response []byte
		fmt.Println("reading…")
		_, response, err := conn.Read(ctx)
		if err != nil {
			fmt.Println("Error reading response:", err)
		}
		outputMutex.Lock()
		fmt.Println("Received: ", string(response))
		outputMutex.Unlock()
	}
}
