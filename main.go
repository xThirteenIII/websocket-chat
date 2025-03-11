package main

import (
	"errors"
	"log"
	"net"
	"os"
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

    return nil
}
