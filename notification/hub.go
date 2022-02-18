package notification

import (
	"log"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

type Hub struct {
	users      map[int32]*User
	notify     chan db.Notification
	register   chan *User
	unregister chan *User
}

func NewHub() *Hub {
	return &Hub{
		users:      make(map[int32]*User),
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
			h.users[user.userID] = user

		case user := <-h.unregister:
			log.Println("User unregister: ", user)
			if _, ok := h.users[user.userID]; ok {
				delete(h.users, user.userID)
				close(user.send)
			}

		case notification := <-h.notify:
			if _, ok := h.users[notification.UserID]; ok {
				h.users[notification.NotificationID].send <- notification
			}
		}
	}
}

func (h *Hub) Notify(userID int32, notification db.Notification) {
	h.notify <- notification
}

func (h *Hub) Register(user *User) {
	h.register <- user
}

func (h *Hub) Unregister(user *User) {
	h.unregister <- user
}
