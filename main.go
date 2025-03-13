package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
	"websocket-chat/src/chat"
)

//import "github.com/coder/websocket"

/*
   Client (HTTP Request with Upgrade header)
             |
             v
   Server (Responds with 101 Switching Protocols)
             |
             v
    Persistent WebSocket Connection Established
             |
             v
    <--- Bidirectional, frame-based messaging --->
*/

// Documentation @:  https://www.rfc-editor.org/rfc/rfc6455.html#section-1.1
func main(){

    log.SetFlags(0)

    err := run()
    if err != nil {
        log.Fatal(err)
    }
}

// run initializes the chatServer and then
// starts a http.Server for the passed in address.
func run() error {

    // If user does not provide at least 1 argument, return error.
    if len(os.Args) < 2 {
        return errors.New("please provide an address to listen on as the first argument")
    }

    // os.Args[0] is the executable command
    // A Listener is a generic network listener for stream-oriented protocols.
    //
    // Multiple goroutines may invoke methods on a Listener simultaneously.
	// Accept() (Conn, error)
	// Close()  error
	// Addr()   Addr
    listener, err := net.Listen("tcp", os.Args[1])
    if err != nil {
        return err
    }

    log.Printf("listening on ws://%v", listener.Addr())

    cs := chat.NewChatServer()
    server := &http.Server{
        Handler: cs,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // why this error handling manouvre with channels?
    // cant just return the error and check it?
    errCh := make(chan error, 1)
    go func(){
        errCh <- server.Serve(listener)
    }()

    signalsCh := make(chan os.Signal, 1)
    // Notify causes package signal to relay incoming signals to signalsCh.
    signal.Notify(signalsCh, os.Interrupt)
    select {
    case err := <- errCh:
        log.Printf("failed to serve: %v", err)
    case sig := <- signalsCh:
        log.Printf("terminating: %v", sig)
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
    defer cancel()

    return server.Shutdown(ctx)
}
