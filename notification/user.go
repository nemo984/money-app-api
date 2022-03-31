package notification

import (
	"log"

	"github.com/gorilla/websocket"
	db "github.com/nemo984/money-app-api/db/sqlc"
)

type User struct {
	ws *websocket.Conn

	userID int32
	send   chan db.Notification
}

func NewUser(ws *websocket.Conn, userID int32) *User {
	return &User{
		send:   make(chan db.Notification),
		ws:     ws,
		userID: userID,
	}
}

// write writes a message with the given message type and payload.
func (u *User) write(mt int, payload []byte) error {
	return u.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (u *User) listen() {
	log.Printf("%v user is listening for message\n", u)
	defer func() {
		log.Printf("%v user stopped listening for message\n", u)
		u.ws.Close()
	}()
	for {
		select {
		case notification, ok := <-u.send:
			if !ok {
				// The hub closed the channel.
				u.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := u.ws.WriteJSON(notification); err != nil {
				return
			}
		}
	}
}
