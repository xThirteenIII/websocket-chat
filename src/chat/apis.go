package chat

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

// subscribeHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (cs *chatServer) subscribeHandler(w http.ResponseWriter, r *http.Request){
    err := cs.subscribe(w, r)
    if errors.Is(err, context.Canceled) {
        return
    }

    if websocket.CloseStatus(err) == websocket.StatusNormalClosure || 
        websocket.CloseStatus(err) == websocket.StatusGoingAway {
            return
    }
    if err != nil {
        cs.logf("error subscribing: %v", err)
        return
    }


}

func (cs *chatServer) subscribe(w http.ResponseWriter, r *http.Request) error {

    var mu sync.Mutex
    var connection *websocket.Conn
    var closed bool
    newSubscriber := &subscriber{
        messagesCh: make(chan []byte, cs.subscriberMessageBuffer),
        closeSlow: func(){
            mu.Lock()
            defer mu.Unlock()
            closed = true
            if connection != nil {
                connection.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
            }
        },
    }
    cs.addSubscriber(newSubscriber)
    defer cs.deleteSubscriber(newSubscriber)

    // Accept accepts a WebSocket handshake from a client and upgrades the
    // the connection to a WebSocket.
    //
    // Accept will not allow cross origin requests by default.
    // See the InsecureSkipVerify and OriginPatterns options to allow cross origin requests.
    //
    // Accept will write a response to w on all errors.
    connection2, err := websocket.Accept(w, r, nil)
    if err != nil {
        return fmt.Errorf("error in WebSocket handshake. Here's why: %v", err)
    }
    mu.Lock()
    if closed {
        mu.Unlock()
        return net.ErrClosed
    }
    connection = connection2
    mu.Unlock()
    return nil
}

func (cs *chatServer) addSubscriber(s *subscriber) {
    
}

func (cs *chatServer) deleteSubscriber(s *subscriber) {
    
}

func (cs *chatServer) publishHandler(w http.ResponseWriter, r *http.Request) {
    
}
