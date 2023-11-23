package bucketo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thepabloaguilar/bucketo"
)

func TestInMemoryTokenStorage_Success(t *testing.T) {
	t.Parallel()

	// SETUP
	storage := bucketo.NewInMemoryTokenStorage(0)
	tokensToSet := []int64{100, 50, 42, 10000}

	// TEST/ASSERT
	for _, tokens := range tokensToSet {
		err := storage.SetTokens(tokens)
		require.NoError(t, err)

		actualTokens, err := storage.GetTokens()
		require.NoError(t, err)
		require.Equal(t, tokens, actualTokens)
	}
}
