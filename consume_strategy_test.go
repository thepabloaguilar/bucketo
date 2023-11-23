package bucketo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thepabloaguilar/bucketo"
)

func TestStaticConsume_Success(t *testing.T) {
	t.Parallel()

	// SETUP
	expectedTokens := int64(1)
	strategy, err := bucketo.NewStaticConsume(expectedTokens)
	require.NoError(t, err)

	// TEST
	actualTokens, err := strategy(nil)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, expectedTokens, actualTokens)
}

func TestStaticConsume_Failure(t *testing.T) {
	t.Parallel()

	// SETUP
	expectedErr := bucketo.ErrTokensToConsumeNegative

	// TEST
	_, err := bucketo.NewStaticConsume(-1)

	// ASSERT
	require.ErrorIs(t, err, expectedErr)
}

func TestNewDynamicConsume_Success(t *testing.T) {
	t.Parallel()

	// SETUP
	strategy := bucketo.NewDynamicConsume()
	expectedTokens := int64(42)

	// TEST
	actualTokens, err := strategy(expectedTokens)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, expectedTokens, actualTokens)
}

func TestNewDynamicConsume_Failure(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		args        any
		expectedErr error
	}{
		{
			name:        "should return an error when argument is not integer",
			args:        "not an integer",
			expectedErr: bucketo.ErrConsumeArgumentIsNotInteger,
		},
		{
			name:        "should return an error when argument is less than zero",
			args:        int64(-1),
			expectedErr: bucketo.ErrTokensToConsumeNegative,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// SETUP
			strategy := bucketo.NewDynamicConsume()

			// TEST
			_, err := strategy(tc.args)

			// ASSERT
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestExpressionsConsume_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		expressions    []bucketo.ConsumeExpression
		args           map[string]any
		expectedTokens int64
	}{
		{
			name: "should return the value from the first expression correctly",
			expressions: []bucketo.ConsumeExpression{
				{
					Expression: "arg > 10",
					Tokens:     10,
				},
				{
					Expression: "arg > 5",
					Tokens:     30,
				},
			},
			args: map[string]any{
				"arg": 20,
			},
			expectedTokens: 10,
		},
		{
			name: "should return the default value correctly",
			expressions: []bucketo.ConsumeExpression{
				{
					Expression: "arg > 10",
					Tokens:     10,
				},
				{
					Expression: "true",
					Tokens:     1,
				},
			},
			args: map[string]any{
				"arg": 10,
			},
			expectedTokens: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// SETUP
			strategy, err := bucketo.NewExpressionsConsume(tc.expressions)
			require.NoError(t, err)

			// TEST
			actualTokens, err := strategy(tc.args)

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, tc.expectedTokens, actualTokens)
		})
	}
}

func TestExpressionsConsume_Failure(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when there are negative token numbers", func(t *testing.T) {
		// SETUP
		expectedEr := bucketo.ErrTokensToConsumeNegative

		// TEST
		_, err := bucketo.NewExpressionsConsume([]bucketo.ConsumeExpression{
			{
				Expression: "true",
				Tokens:     -1,
			},
		})

		// ASSERT
		require.ErrorIs(t, err, expectedEr)
	})

	t.Run("should return an error when none of the expressions return true", func(t *testing.T) {
		// SETUP
		expectedEr := bucketo.ErrNoExpressionMatched

		strategy, err := bucketo.NewExpressionsConsume([]bucketo.ConsumeExpression{
			{
				Expression: "false",
				Tokens:     1,
			},
		})
		require.NoError(t, err)

		// TEST
		_, err = strategy(nil)

		// ASSERT
		require.ErrorIs(t, err, expectedEr)
	})
}
