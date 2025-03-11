package chat

import (
	"log"
	"net/http"
	"sync"

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

