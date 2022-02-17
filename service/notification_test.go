package service

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

func TestCreateNotifications(t *testing.T) {
	testCases := []db.CreateNotificationParams{
		{
			UserID: 1,
			Description: sql.NullString{
				String: util.RandomString(20),
				Valid:  true,
			},
			Priority: db.NotificationPriorityLow,
			Type:     "Urgent",
		},
		{
			UserID: 1,
			Description: sql.NullString{
				String: util.RandomString(20),
				Valid:  true,
			},
			Priority: db.NotificationPriorityMedium,
			Type:     "Cool",
		},
		{
			UserID: 1,
			Description: sql.NullString{
				String: util.RandomString(20),
				Valid:  true,
			},
			Priority: db.NotificationPriorityHigh,
			Type:     "Your Mom!",
		},
	}

	for _, tc := range testCases {
		category, err := testService.CreateNotification(context.Background(), tc)
		require.NoError(t, err)
		require.NotEmpty(t, category)
	}

}
