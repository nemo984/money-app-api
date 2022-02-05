package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomBudget(t *testing.T, userID int32) Budget {
	arg := CreateBudgetParams{
		UserID:       userID,
		CategoryName: randomCategory(),
		Amount:       fmt.Sprint(util.RandomInt(1000, 20000)),
		EndDate: sql.NullTime{
			Time:  time.Now().AddDate(0, 1, 0),
			Valid: true,
		},
	}

	budget, err := testQueries.CreateBudget(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, budget)

	require.Equal(t, arg.UserID, budget.UserID)
	require.Equal(t, arg.CategoryName, budget.CategoryName)
	requireAmountEqual(t, arg.Amount, budget.Amount)
	require.WithinDuration(t, arg.EndDate.Time, budget.EndDate.Time, time.Second)

	require.NotZero(t, budget.CreatedAt)

	return budget
}

func TestCreateBudget(t *testing.T) {
	user := createRandomUser(t)
	createRandomBudget(t, user.UserID)
}

func TestDeleteBudget(t *testing.T) {
	user := createRandomUser(t)
	i := createRandomBudget(t, user.UserID)
	err := testQueries.DeleteBudget(context.Background(), i.UserID)
	require.NoError(t, err)
}
