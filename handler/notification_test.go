package handler_test

import (
	"database/sql"
	"os"
	"testing"

	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/handler"
	"github.com/nemo984/money-app-api/notification"
	"github.com/nemo984/money-app-api/util"
)

var hub *notification.Hub

func TestMain(m *testing.M) {
	hub = notification.NewHub()
	server := handler.NewServer(nil, hub)
	go hub.Run()
	go server.Start(":8081")
	os.Exit(m.Run())
}

func TestSendNotification(t *testing.T) {
	var userID int32 = 1
	hub.Notify(userID, db.Notification{
		NotificationID: int32(util.RandomInt(0, 100)),
		UserID:         userID,
		Description: sql.NullString{
			String: util.RandomString(100),
			Valid:  true,
		},
		Type:     "Budgets",
		Priority: db.NotificationPriorityHigh,
		Read:     false,
	})
}
