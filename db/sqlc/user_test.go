package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomString(10),
		Name: sql.NullString{String: util.RandomString(10), Valid: true},
		Password: util.RandomString(10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)
	
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)	
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.UserID, user2.UserID)
	require.Equal(t, user.Password, user2.Password)
	require.Equal(t, user.Name, user2.Name)
	require.WithinDuration(t, user.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserParams{
		UserID: user.UserID,
		Username: util.RandomString(10),
		Name: sql.NullString{String: util.RandomString(10), Valid: true},
		Password: util.RandomString(10),
		ProfileUrl: sql.NullString{String: util.RandomString(30), Valid: true},
	}
	
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.UserID, user2.UserID)
	require.WithinDuration(t, user.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
	
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Name, user2.Name)
	require.Equal(t, arg.Password, user2.Password)
	require.Equal(t, arg.ProfileUrl, user2.ProfileUrl)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}