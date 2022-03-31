package notification

import (
	"log"
	"sync"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

type Hub struct {
	usersMap   sync.Map
	notify     chan db.Notification
	register   chan *User
	unregister chan *User
}

func New() *Hub {
	return &Hub{
		notify:     make(chan db.Notification),
		register:   make(chan *User),
		unregister: make(chan *User),
	}
}

func (h *Hub) Run() {
	log.Println("Hub is Listening")
	defer log.Println("Hub is dead")
	for {
		select {
		case user := <-h.register:
			log.Println("New User connection: ", user)
			h.usersMap.Store(user.userID, user)

		case user := <-h.unregister:
			log.Println("User unregister: ", user)
			if _, ok := h.usersMap.LoadAndDelete(user.userID); ok {
				close(user.send)
			}

		case notification := <-h.notify:
			if u, ok := h.usersMap.Load(notification.UserID); ok {
				u.(*User).send <- notification
			}
		}
	}
}

func (h *Hub) Notify(userID int32, notification db.Notification) {
	h.notify <- notification
}

func (h *Hub) Register(user *User) {
	h.register <- user
	user.listen()
}

func (h *Hub) Unregister(user *User) {
	h.unregister <- user
}
