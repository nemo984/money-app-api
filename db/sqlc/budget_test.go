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
func TestGetBudgets(t *testing.T) {
	user := createRandomUser(t)
	var budgets []Budget
	n := 5
	for i := 0; i < n; i++ {
		budget := createRandomBudget(t, user.UserID)
		budgets = append(budgets, budget)
	}
	b2, err := testQueries.GetBudgets(context.Background(), user.UserID)
	require.NoError(t, err)
	require.Equal(t, n, len(b2))
	for i := 0; i < n; i++ {
		require.Equal(t, budgets[i].UserID, b2[i].UserID)
		require.Equal(t, budgets[i].CategoryName, b2[i].CategoryName)
		requireAmountEqual(t, budgets[i].Amount, b2[i].Amount)
	}
}
func TestDeleteBudget(t *testing.T) {
	user := createRandomUser(t)
	i := createRandomBudget(t, user.UserID)
	err := testQueries.DeleteBudget(context.Background(), i.UserID)
	require.NoError(t, err)
}
