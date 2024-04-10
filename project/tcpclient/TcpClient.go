package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("tcpclient", "Create a TCP client to connect to a middleware node")
	nodeIp := parser.String("i", "node-ip", nil)
	nodePort := parser.Int("", "node-port", nil)
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
	nodeAddress := fmt.Sprintf("%v:%v", nodeIp, nodePort)
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		return err
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		message, _ := reader.ReadString('\n')

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		response := make([]byte, 1024)
		_, err = conn.Read(response)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Received:", string(response))
	}
	return conn.Close()
}

func runTests(nodeIp string, nodePort int, filePath string) error {

	nodeAddress := fmt.Sprintf("%v:%v", nodeIp, nodePort)
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		return err
	} else {
		fmt.Println("Connection successful")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	for i := 0; i < 1000; i++ {
		message := "test\n"
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		startTime := time.Now()

		response := make([]byte, 1024)
		_, err = conn.Read(response)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		file.WriteString(latency.String() + "\n")

		fmt.Println("Received:", string(response))
	}
	return conn.Close()
}
