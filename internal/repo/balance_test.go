package repo

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetBalance(t *testing.T) {
	arg := 27
	balance, err := testBalance.GetBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, balance)
}
