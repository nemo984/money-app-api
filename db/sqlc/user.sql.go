// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username, name, password, profile_url
) VALUES (
  $1, $2, $3, $4
)
RETURNING user_id, username, name, password, profile_url, created_at
`

type CreateUserParams struct {
	Username   string         `json:"username"`
	Name       sql.NullString `json:"name"`
	Password   string         `json:"password"`
	ProfileUrl sql.NullString `json:"profile_url"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Name,
		arg.Password,
		arg.ProfileUrl,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.ProfileUrl,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, userID)
	return err
}

const getUser = `-- name: GetUser :one
SELECT user_id, username, name, password, profile_url, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.ProfileUrl,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT user_id, username, name, password, profile_url, created_at FROM users
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.ProfileUrl,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET 
    username = CASE WHEN $1::boolean
        THEN $2::VARCHAR(20) ELSE username END,

    name = CASE WHEN $3::boolean
        THEN $4::VARCHAR(20) ELSE name END,

    password = CASE WHEN $5::boolean
        THEN $6::VARCHAR ELSE password END,

    profile_url = CASE WHEN $7::boolean
        THEN $8::VARCHAR ELSE profile_url END

WHERE user_id = $9
RETURNING user_id, username, name, password, profile_url, created_at
`

type UpdateUserParams struct {
	UsernameDoUpdate bool   `json:"username_do_update"`
	Username         string `json:"username"`
	NameDoUpdate     bool   `json:"name_do_update"`
	Name             string `json:"name"`
	PasswordDoUpdate bool   `json:"password_do_update"`
	Password         string `json:"password"`
	ProfileDoUpdate  bool   `json:"profile_do_update"`
	ProfileUrl       string `json:"profile_url"`
	UserID           int32  `json:"user_id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.UsernameDoUpdate,
		arg.Username,
		arg.NameDoUpdate,
		arg.Name,
		arg.PasswordDoUpdate,
		arg.Password,
		arg.ProfileDoUpdate,
		arg.ProfileUrl,
		arg.UserID,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.ProfileUrl,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPicture = `-- name: UpdateUserPicture :one
UPDATE users
SET profile_url = $2
WHERE user_id = $1
RETURNING user_id, username, name, password, profile_url, created_at
`

type UpdateUserPictureParams struct {
	UserID     int32          `json:"user_id"`
	ProfileUrl sql.NullString `json:"profile_url"`
}

func (q *Queries) UpdateUserPicture(ctx context.Context, arg UpdateUserPictureParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPicture, arg.UserID, arg.ProfileUrl)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.ProfileUrl,
		&i.CreatedAt,
	)
	return i, err
}
