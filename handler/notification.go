package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/service"
)

func (s *Server) getNotifications(c *gin.Context) {
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
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
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
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
	userPayload := c.MustGet(authorizationPayload).(service.JWTClaims)
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
