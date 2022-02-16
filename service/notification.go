package service

import (
	"context"

	db "github.com/nemo984/money-app-api/db/sqlc"
)

func (s *service) CreateNotification(ctx context.Context, args db.CreateNotificationParams) (db.Notification, error) {
	noti, err := s.db.CreateNotification(ctx, args)
	if err != nil {
		return db.Notification{}, err
	}

	return noti, nil
}

func (s *service) UpdateNotification(ctx context.Context, args db.UpdateNotificationParams) (db.Notification, error) {
	noti, err := s.db.UpdateNotification(ctx, args)
	if err != nil {
		return db.Notification{}, err
	}

	return noti, nil
}
