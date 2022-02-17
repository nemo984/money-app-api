package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/notification"
	"github.com/nemo984/money-app-api/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) wsNotificationHandler(hub *notification.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := s.service.VerifyToken(c, c.Query("token"))
		if err != nil {
			err := errors.New("missing or invalid jwt token")
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := notification.NewUser(ws, claims.UserID)
		hub.Register(user)
		user.Listen()
	}
}

func (s *Server) getNotifications(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	notifications, err := s.service.GetNotifications(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notifications)
}

type notificationURI struct {
	NotificationID int32 `uri:"id"`
}

func (s *Server) updateNotification(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req notificationURI
	if err := c.ShouldBindUri(&req); err != nil {
		handleError(c, err)
		return
	}

	notification, err := s.service.UpdateNotification(c, userPayload.UserID, db.UpdateNotificationParams{
		NotificationID: req.NotificationID,
		Read:           true,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notification)
}

func (s *Server) updateAllNotifications(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	args := db.UpdateNotificationsParams{
		UserID: userPayload.UserID,
		Read:   true,
	}

	notifications, err := s.service.UpdateNotifications(c, args)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notifications)
}
