package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomNotification(t *testing.T, userID int32) Notification {
	arg := CreateNotificationParams{
		UserID: userID,
		Description: sql.NullString{
			String: util.RandomString(100),
			Valid:  true,
		},
		Type:     util.RandomString(15),
		Priority: NotificationPriorityHigh,
	}

	n, err := testQueries.CreateNotification(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, n)

	require.Equal(t, arg.UserID, n.UserID)
	require.Equal(t, arg.Description.String, n.Description.String)
	require.Equal(t, arg.Type, n.Type)
	require.Equal(t, arg.Priority, n.Priority)
	require.NotZero(t, n.CreatedAt)

	return n
}

func TestCreateNotification(t *testing.T) {
	user := createRandomUser(t)
	createRandomNotification(t, user.UserID)
}

func TestGetNotifications(t *testing.T) {
	user := createRandomUser(t)
	var notifications []Notification
	n := 5
	for i := 0; i < n; i++ {
		notification := createRandomNotification(t, user.UserID)
		notifications = append(notifications, notification)
	}
	n2, err := testQueries.GetNotifications(context.Background(), user.UserID)
	require.NoError(t, err)
	require.Equal(t, n, len(n2))
	for i := 0; i < n; i++ {
		require.Equal(t, notifications[i].UserID, n2[i].UserID)
		require.Equal(t, notifications[i].Description.String, n2[i].Description.String)
		require.Equal(t, notifications[i].Type, n2[i].Type)
		require.Equal(t, notifications[i].Priority, n2[i].Priority)
	}
}

func TestUpdateNotification(t *testing.T) {
	user := createRandomUser(t)
	n := createRandomNotification(t, user.UserID)
	arg := UpdateNotificationParams{
		NotificationID: n.NotificationID,
		Read: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}
	n2, err := testQueries.UpdateNotification(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, user.UserID, n2.UserID)
	require.Equal(t, arg.Read.Bool, n2.Read.Bool)
}

func TestDeleteNotification(t *testing.T) {
	user := createRandomUser(t)
	i := createRandomNotification(t, user.UserID)
	err := testQueries.DeleteNotification(context.Background(), i.UserID)
	require.NoError(t, err)
}
