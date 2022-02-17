package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) GetNotifications(ctx context.Context, userID int32) ([]db.Notification, error) {
	notifications, err := s.db.GetNotifications(ctx, userID)
	if err != nil {
		return []db.Notification{}, nil
	}

	return notifications, nil
}

func (s *service) CreateNotification(ctx context.Context, args db.CreateNotificationParams) (db.Notification, error) {
	noti, err := s.db.CreateNotification(ctx, args)
	if err != nil {
		return db.Notification{}, err
	}

	return noti, nil
}

func (s *service) UpdateNotification(ctx context.Context, userID int32, args db.UpdateNotificationParams) (db.Notification, error) {
	notification, err := s.db.GetNotification(ctx, args.NotificationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Notification{}, AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("no notification with that id"),
			}
		}
		return db.Notification{}, fmt.Errorf("db get notification error: %v", err)
	}

	if notification.UserID != userID {
		return db.Notification{}, AppError{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("not your expense"),
		}
	}

	noti, err := s.db.UpdateNotification(ctx, args)
	if err != nil {
		return db.Notification{}, err
	}

	return noti, nil
}

func (s *service) UpdateNotifications(ctx context.Context, args db.UpdateNotificationsParams) ([]db.Notification, error) {
	notifications, err := s.db.UpdateNotifications(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return []db.Notification{}, AppError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("no notification with that id"),
			}
		}
		return []db.Notification{}, err
	}

	return notifications, err
}
