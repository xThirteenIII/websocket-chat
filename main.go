package main

import "golang.org/x/net/websocket"

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

func main(){
}
