package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"project/cloud/httpserver"
	"project/cloud/socketserver"
	"time"

	"github.com/akamensky/argparse"
)

func main() {
	parser := *argparse.NewParser("server", "Create a new server instance")
	serverType := parser.String("t", "type", nil)
	serverAddress := parser.String("s", "server-address", nil)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if *serverType == "websocket" {
		fmt.Println("Let's make a websocket")
		err = StartWebsocketServer(*serverAddress)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else if *serverType == "socket" {
		err = StartSocketServer(*serverAddress)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else if *serverType == "api" {
		fmt.Println("Not yet implemented")
	} else {
		fmt.Println("Please set server type to websocket, socket, or api")
	}
}

func StartWebsocketServer(serverAddress string) error {
	ctx := context.WithoutCancel(context.Background())

	l, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	httpServer := httpserver.CreateHttpCloudServer()
	server := &http.Server{
		Handler:      httpServer,
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

func StartSocketServer(serverAddress string) error {
	server := socketserver.CreateSocketServer()
	l, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	for {
		clientConn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
		}

		go server.HandleTcp(clientConn)

	}
}
