package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	
	arg := CreateUserParams{
		Username: "Jack",
		Name: sql.NullString{String: "Jack' O reilly", Valid: true},
		Password: "2420",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)

}