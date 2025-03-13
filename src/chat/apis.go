package chat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

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
    defer connection.CloseNow()

    ctx := connection.CloseRead(context.Background())

    for {
        select { 
        case msg := <- newSubscriber.messagesCh:
            err := writeTimeout(ctx, time.Second * 5, connection, msg)
            if err != nil {
                return err
            }

        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// addSubscriber registers a subscriber.
func (cs *chatServer) addSubscriber(s *subscriber) {
    cs.subscribersMu.Lock()
    cs.subscribers[s] = struct{}{}
    cs.subscribersMu.Unlock()
    
}

// addSubscriber deletes the given subscriber.
func (cs *chatServer) deleteSubscriber(s *subscriber) {
    cs.subscribersMu.Lock()
    delete(cs.subscribers, s)
    cs.subscribersMu.Unlock()
}

// publishHandler reads the request body with a limit of 8192 bytes and then 
// publishes the received message.
func (cs *chatServer) publishHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
        return
    }
    body := http.MaxBytesReader(w, r.Body, 8192)
    msg, err := io.ReadAll(body)
    if err != nil {
        http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
        return
    }

    cs.publish(msg)

    w.WriteHeader(http.StatusAccepted)
    
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subs are blocked.
func (cs *chatServer) publish(msg []byte){
    cs.subscribersMu.Lock()
    defer cs.subscribersMu.Unlock()

    cs.publishLimiter.Wait(context.Background())

    for s := range cs.subscribers {
        select {
        case s.messagesCh <- msg:
        default:
            go s.closeSlow()

        }
    }

}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
