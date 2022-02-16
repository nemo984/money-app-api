package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{name: "Transportation"},
		{name: "Food"},
		{name: "Shopping"},
	}

	for _, tc := range testCases {
		category, err := testService.CreateCategory(context.Background(), tc.name)
		require.NoError(t, err)
		require.NotEmpty(t, category)
		require.Equal(t, tc.name, category.Name)
	}

}
