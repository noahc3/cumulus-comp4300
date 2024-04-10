package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	parser := argparse.NewParser("wsclient", "Create a client that connects to middlware via websocket")
	nodeIp := parser.String("i", "node-ip", nil)
	nodePort := parser.Int("p", "node-port", nil)
	testFile := parser.String("t", "test", nil)
	parser.Parse(os.Args)

	var err error
	if *testFile != "" {
		err = runTests(*nodeIp, *nodePort, *testFile)
	} else {
		err = run(*nodeIp, *nodePort)
	}

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Program ended successfully")
}

func run(nodeIp string, nodePort int) error {
	// connect to the websocket node using the supplied address
	ctx := context.WithoutCancel(context.Background())

	nodeAddrString := fmt.Sprintf("ws://%v:%v/", nodeIp, nodePort)
	conn, _, err := websocket.Dial(ctx, nodeAddrString, nil)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	// wait for user-inputted commands…
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		message, _ := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		err = wsjson.Write(ctx, conn, message)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println(message)
	}

	return conn.Close(websocket.StatusNormalClosure, "")
}

func runTests(nodeIp string, nodePort int, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// connect to the websocket node using the supplied address
	ctx := context.WithoutCancel(context.Background())

	nodeAddrString := fmt.Sprintf("ws://%v:%v/", nodeIp, nodePort)
	conn, _, err := websocket.Dial(ctx, nodeAddrString, nil)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	for i := 0; i < 1000; i++ {
		// send a test message and measure the latency!
		fmt.Println("Sending: test")
		err = wsjson.Write(ctx, conn, "test\n")
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		startTime := time.Now()

		var message interface{}
		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			fmt.Println("Error:", err)
		}
		endTime := time.Now()

		latency := endTime.Sub(startTime)
		// converts to milliseconds
		file.WriteString(fmt.Sprintf("%v\n", float64(latency.Nanoseconds())/float64(10e6)))
		fmt.Println("Received:", message)
	}

	return conn.Close(websocket.StatusNormalClosure, "")
}
