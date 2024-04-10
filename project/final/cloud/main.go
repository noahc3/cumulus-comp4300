package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"project/final/cloud/httpserver"
	"time"

	"github.com/akamensky/argparse"
)

func main() {
	parser := *argparse.NewParser("server", "Create a new server instance")
	host := parser.String("H", "host", &argparse.Options{Required: false, Help: "Host to listen on", Default: "0.0.0.0"})
	port := parser.String("P", "port", &argparse.Options{Required: false, Help: "Port to listen on", Default: "1111"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
	}

	httpserver.StartLoop()

	err = StartWebsocketServer(*host, *port)
	if err != nil {
		log.Fatal(err)
	}
}

func StartWebsocketServer(host string, port string) error {
	ctx := context.WithoutCancel(context.Background())

	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
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
