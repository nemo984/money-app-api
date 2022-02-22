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
func (s *Server) wsNotificationHandler(c *gin.Context) {
	var query notificationWSToken
	if err := c.ShouldBindQuery(&query); err != nil {
		handleValidationError(c, &query, err)
		return
	}

	claims, err := s.service.VerifyToken(c, query.Token)
	if err != nil {
		err := errors.New("invalid jwt token")
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
	s.hub.Register(user)
	defer s.hub.Unregister(user)
	user.Listen()
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
func (s *Server) getNotifications(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	notifications, err := s.service.GetNotifications(c, userPayload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notifications)
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
func (s *Server) updateNotification(c *gin.Context) {
	userPayload := c.MustGet(AuthorizationPayload).(service.JWTClaims)
	var req notificationURI
	if err := c.ShouldBindUri(&req); err != nil {
		handleValidationError(c, &req, err)
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

// swagger:route PATCH /me/notifications Notifications updateNotifications
// Set all user's notifications to read and returns them
//
// Security:
//  bearerAuth:
//  cookieAuth:
//
// responses:
//  200: notificationsResponse
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
