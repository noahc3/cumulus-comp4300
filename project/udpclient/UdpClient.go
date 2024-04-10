package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("udpclient", "Create a UDP client to connect to a middleware node")
	nodeIp := parser.String("i", "node-ip", nil)
	nodePort := parser.Int("", "node-port", nil)
	clientPort := parser.Int("c", "client-port", nil)
	testFile := parser.String("t", "test", nil)
	parser.Parse(os.Args)

	var err error
	if *testFile != "" {
		err = runTests(*nodeIp, *nodePort, *clientPort, *testFile)
	} else {
		err = run(*nodeIp, *nodePort, *clientPort)
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Program ended successfully")
}

func run(nodeIp string, nodePort int, clientPort int) error {
	udpAddress := &net.UDPAddr{
		IP:   net.ParseIP(nodeIp),
		Port: nodePort,
	}

	udpReturnAddress := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: clientPort,
	}

	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	listenerConn, err := net.ListenUDP("udp", udpReturnAddress)
	if err != nil {
		return err
	}

	for {
		reader := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter text: ")
		reader.Scan()
		message := reader.Text()

		msgJson := fmt.Sprintf("{\"message\": \"%v\", \"returnAddr\": %v}", message, clientPort)

		_, err = conn.Write([]byte(msgJson))
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		response := make([]byte, 1024)
		_, _, err = listenerConn.ReadFromUDP(response)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Received:", string(response))

	}

	return conn.Close()
}

func runTests(nodeIp string, nodePort int, clientPort int, filePath string) error {
	udpAddress := &net.UDPAddr{
		IP:   net.ParseIP(nodeIp),
		Port: nodePort,
	}

	udpReturnAddress := &net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: clientPort,
	}

	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	listenerConn, err := net.ListenUDP("udp", udpReturnAddress)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	for i := 0; i < 1000; i++ {

		msgJson := make(map[string]string)
		msgJson["command"] = "input"
		msgJson["value"] = "test\n"
		msgJson["returnAddr"] = fmt.Sprintf("%v", udpReturnAddress.Port)
		jsonString, err := json.Marshal(msgJson)
		if err != nil {
			return err
		}

		fmt.Println("Writing:", string(jsonString))

		_, err = conn.Write(jsonString)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		startTime := time.Now()

		response := make([]byte, 1024)
		_, _, err = listenerConn.ReadFromUDP(response)
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
