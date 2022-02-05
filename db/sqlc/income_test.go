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

func TestDeleteIncome(t *testing.T) {
	user := createRandomUser(t)
	i := createRandomIncome(t, user.UserID)
	err := testQueries.DeleteExpense(context.Background(), i.UserID)
	require.NoError(t, err)
}
