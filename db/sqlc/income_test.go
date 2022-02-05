package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/nemo984/money-app-api/util"
	"github.com/stretchr/testify/require"
)

//TODO:
func randomIncomeType() string {
	return "Passive"
}

func createRandomIncome(t *testing.T, userID int32) Income {
	arg := CreateIncomeParams{
		UserID:         userID,
		IncomeTypeName: randomIncomeType(),
		Description: sql.NullString{
			String: util.RandomString(50),
			Valid:  true,
		},
		Amount:    fmt.Sprint(util.RandomInt(20000, 40000)),
		Frequency: DateFrequencyMonth,
	}

	income, err := testQueries.CreateIncome(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, income)

	require.Equal(t, arg.UserID, income.UserID)
	require.Equal(t, arg.IncomeTypeName, income.IncomeTypeName)
	require.Equal(t, arg.Description, income.Description)
	require.Equal(t, arg.Frequency, income.Frequency)

	require.NotZero(t, income.CreatedAt)

	return income
}

func TestCreateIncome(t *testing.T) {
	user := createRandomUser(t)
	createRandomIncome(t, user.UserID)
}

func TestGetIncomes(t *testing.T) {
	user := createRandomUser(t)
	var incomes []Income
	n := 5
	for i := 0; i < n; i++ {
		income := createRandomIncome(t, user.UserID)
		incomes = append(incomes, income)
	}
	i2, err := testQueries.GetIncomes(context.Background(), user.UserID)
	require.NoError(t, err)
	require.Equal(t, n, len(i2))
	for i := 0; i < n; i++ {
		require.Equal(t, incomes[i].UserID, i2[i].UserID)
		require.Equal(t, incomes[i].IncomeTypeName, i2[i].IncomeTypeName)
		require.Equal(t, incomes[i].Description, i2[i].Description)
		require.Equal(t, incomes[i].Frequency, i2[i].Frequency)
	}
}

func TestDeleteIncome(t *testing.T) {
	user := createRandomUser(t)
	i := createRandomIncome(t, user.UserID)
	err := testQueries.DeleteExpense(context.Background(), i.UserID)
	require.NoError(t, err)
}
