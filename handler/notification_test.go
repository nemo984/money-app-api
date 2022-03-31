package handler_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/handler"
	"github.com/nemo984/money-app-api/notification"
	"github.com/nemo984/money-app-api/service"
	mockservice "github.com/nemo984/money-app-api/service/mock"
	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

var (
	hub       *notification.Hub
	URL       = "ws://127.0.0.1:8081/api/notifications-ws"
	testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOm51bGwsInVzZXJfaWQiOjF9.5_ws3YIck_0n4ayTwO_ufsZJn-nXFk0z5BR76v3UN9A"
)

func setUpServer(t *testing.T, userID int32) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mockservice.NewMockService(mockCtrl)
	mockService.
		EXPECT().
		VerifyToken(gomock.Any(), gomock.Any()).
		Return(service.JWTClaims{
			MapClaims: make(jwt.MapClaims),
			UserID:    userID,
		}, nil).
		AnyTimes()

	hub = notification.New()
	s := handler.New(mockService, hub)
	go hub.Run()
	go s.Start(":8081")
}

func TestSendNotification(t *testing.T) {
	var userID int32 = 1
	setUpServer(t, userID)

	N := 3
	url := URL + "?token=" + testToken
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	defer ws.Close()

	for i := 0; i < N; i++ {
		sendNoti := db.Notification{
			NotificationID: int32(util.RandomInt(0, 100)),
			UserID:         userID,
			Description: sql.NullString{
				String: util.RandomString(100),
				Valid:  true,
			},
			Type:     "Budgets",
			Priority: db.NotificationPriorityHigh,
			Read:     false,
		}
		hub.Notify(userID, sendNoti)

		_, p, err := ws.ReadMessage()
		require.NoError(t, err)
		var noti db.Notification
		err = json.Unmarshal(p, &noti)
		require.NoError(t, err)
		require.NotEmpty(t, noti)
		require.Equal(t, sendNoti.NotificationID, noti.NotificationID)
		require.Equal(t, sendNoti.UserID, noti.UserID)
		require.Equal(t, sendNoti.Description, noti.Description)
		require.Equal(t, sendNoti.Type, noti.Type)
		require.Equal(t, sendNoti.Priority, noti.Priority)
		require.Equal(t, sendNoti.Read, noti.Read)
	}
}
