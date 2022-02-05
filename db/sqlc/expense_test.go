package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

//TODO: .
func randomCategory() string {
	return "Transportation"
}

var date = [...]DateFrequency{DateFrequencyDay, DateFrequencyMonth, DateFrequencyWeek, DateFrequencyYear}

func randomFrequency() DateFrequency {
	return date[util.RandomInt(0, len(date)-1)]
}

func createRandomExpense(t *testing.T, userID int32) Expense {
	arg := CreateExpenseParams{
		UserID:       userID,
		CategoryName: randomCategory(),
		Amount:       fmt.Sprint(util.RandomInt(100, 20000)),
		Frequency:    randomFrequency(),
		Note: sql.NullString{
			String: util.RandomString(50),
			Valid:  true,
		},
	}

	expense, err := testQueries.CreateExpense(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, expense)
	requireAmountEqual(t, arg.Amount, expense.Amount)
	require.Equal(t, arg.UserID, expense.UserID)
	require.Equal(t, arg.CategoryName, expense.CategoryName)

	require.Equal(t, arg.Frequency, expense.Frequency)
	require.Equal(t, arg.Note, expense.Note)

	require.NotZero(t, expense.ExpenseID)
	require.NotZero(t, expense.CreatedAt)

	return expense
}

func TestCreateExpense(t *testing.T) {
	user := createRandomUser(t)
	createRandomExpense(t, user.UserID)
}

func TestGetExpenses(t *testing.T) {
	user := createRandomUser(t)
	var expenses []Expense
	n := 5
	for i := 0; i < n; i++ {
		expense := createRandomExpense(t, user.UserID)
		expenses = append(expenses, expense)
	}
	e2, err := testQueries.GetExpenses(context.Background(), user.UserID)
	require.NoError(t, err)
	require.Equal(t, n, len(e2))
	for i := 0; i < n; i++ {
		require.Equal(t, expenses[i].UserID, e2[i].UserID)
		require.Equal(t, expenses[i].CategoryName, e2[i].CategoryName)
		require.Equal(t, expenses[i].Frequency, e2[i].Frequency)
		require.Equal(t, expenses[i].Note, e2[i].Note)
	}
}

func TestDeleteExpense(t *testing.T) {
	user := createRandomUser(t)
	e := createRandomExpense(t, user.UserID)
	err := testQueries.DeleteExpense(context.Background(), e.ExpenseID)
	require.NoError(t, err)
}

func requireAmountEqual(t *testing.T, actual, expected string) {
	argAmount, err := strconv.ParseFloat(actual, 64)
	require.NoError(t, err)
	expAmount, err := strconv.ParseFloat(expected, 64)
	require.NoError(t, err)
	require.Equal(t, int(argAmount), int(expAmount))
}
