package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateIncomeType(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{name: "Passive"},
		{name: "Full-Time Job"},
		{name: "Part-Time Job"},
	}

	for _, tc := range testCases {
		incomeType, err := testService.CreateIncomeType(context.Background(), tc.name)
		require.NoError(t, err)
		require.NotEmpty(t, incomeType)
		require.Equal(t, tc.name, incomeType.Name)
	}

}
