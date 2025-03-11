package chat

import (
	"context"
	"errors"
	"net/http"

	"github.com/coder/websocket"
)

// subscribeHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (cs *chatServer) subscriberHandler(w http.ResponseWriter, r *http.Request){
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
