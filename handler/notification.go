package handler

import (
	"fmt"
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

// swagger:parameters listenNotifications
type notificationWSToken struct {
	// Auth token for the user
	// required: true
	Token string `json:"token" form:"token"`
}

// swagger:route GET /notifications-ws Notifications listenNotifications
// Listen for notifications
// Schemes: ws
// responses:
//  200: notificationResponse
func (h *handler) WSNotificationHandler(c *gin.Context) (interface{}, int, error) {
	var query notificationWSToken
	if err := c.ShouldBindQuery(&query); err != nil {
		err := validationErrors(&query, err)
		return nil, 0, err
	}

	claims, err := h.service.VerifyToken(c, query.Token)
	if err != nil {
		return nil, 0, err
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		err := fmt.Errorf("handler upgrade error: %w", err)
		return nil, 0, err
	}

	user := notification.NewUser(ws, claims.UserID)
	h.hub.Register(user)
	defer h.hub.Unregister(user)
	return "", http.StatusNoContent, nil
}

// A list of notifications
// swagger:response notificationsResponse
type notificationsResponse struct {
	// User's notifications
	// in:body
	Body []db.Notification
}

// Notification
// swagger:response notificationResponse
type notificationResponse struct {
	Body db.Notification
}

// swagger:route GET /me/notifications Notifications listNotifications
// Returns a list of notifications of the user
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: notificationsResponse
func (h *handler) getNotifications(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	notifications, err := h.service.GetNotifications(c, userPayload.UserID)
	if err != nil {
		return nil, 0, err
	}

	return notifications, http.StatusOK, nil
}

// swagger:parameters updateNotification
type notificationURI struct {
	// The id of the notification to update
	// in: path
	// required: true
	// min: 1
	NotificationID int32 `uri:"id" binding:"min=1"`
}

// swagger:route PATCH /me/notifications/{id} Notifications updateNotification
// Set the notification to read and returns the notification
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: notificationResponse
func (h *handler) updateNotification(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req notificationURI
	if err := c.ShouldBindUri(&req); err != nil {
		err := validationErrors(&req, err)
		return nil, 0, err
	}

	notification, err := h.service.UpdateNotification(c, userPayload.UserID, db.UpdateNotificationParams{
		NotificationID: req.NotificationID,
		Read:           true,
	})
	if err != nil {
		return nil, 0, err
	}

	return notification, http.StatusOK, nil
}

// swagger:route PATCH /me/notifications Notifications updateNotifications
// Set all user's notifications to read and returns them
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: notificationsResponse
func (h *handler) updateAllNotifications(c *gin.Context) (interface{}, int, error) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	args := db.UpdateNotificationsParams{
		UserID: userPayload.UserID,
		Read:   true,
	}

	notifications, err := h.service.UpdateNotifications(c, args)
	if err != nil {
		return nil, 0, err
	}

	return notifications, http.StatusOK, nil
}
