package chat

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

//"net/http"

//"github.com/coder/websocket"

// chatServer enables broadcasting to a range of subscribers.
type chatServer struct {

    // subscriberMessageBuffer holds the maximum messages per subscriber that can be queued
    // before it is kicked out.
    //
    // Defaults to 16.
    subscriberMessageBuffer int 

    // publishLimiter limits the publishing to the endpoint.
    //
    // Defaults to one publish every 100ms with a burst of 8.
    publishLimiter *rate.Limiter

    // logf controls where logs are sent.
    // It mirrors log.Printf, log.Fatalf functions arguments,
    // e.g. When creating a new chatServer:
    // cs := chatServer {
    //          ...
    //          logf: log.Fatalf,
    //          ...
    //       }
    //
    // Defaults to log.Printf.
    logf func(f string, v ...any)

    // serveMux routes the various endpoints to the appropiate handler.
    //
    // ServeMux is an HTTP request multiplexer.
    // It matches the URL of each incoming request against a list of registered
    // patterns and calls the handler for the pattern that
    // most closely matches the URL.
    // In general, a pattern looks like
    //
    //	[METHOD ][HOST]/[PATH]
    serveMux    http.ServeMux

    subscribersMu sync.Mutex
    subscribers map[*subscriber]struct{}
}

// Messages are sent on the messages channel and if the client cannot keepup
// with the messages, closeSlow is called.
type subscriber struct {
    messagesCh  chan []byte
    closeSlow   func()
}

func (cs *chatServer)ServeHTTP(w http.ResponseWriter, r *http.Request) {
    cs.serveMux.ServeHTTP(w, r)
}

// newChatServer constructs and returns a new chatServer, filled
// with defaults.
func NewChatServer() *chatServer {

    cs := &chatServer{
        subscriberMessageBuffer:    16,
        logf:                       log.Printf,
        subscribers:                make(map[*subscriber]struct{}),
        publishLimiter:             rate.NewLimiter(rate.Every(time.Millisecond * 100), 8),
    }
    cs.serveMux.Handle("/", http.FileServer(http.Dir(".")))
    cs.serveMux.HandleFunc("/subscribe", cs.subscribeHandler)
    cs.serveMux.HandleFunc("/publish", cs.publishHandler)

    return cs
}
