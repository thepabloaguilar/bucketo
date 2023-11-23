package bucketo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thepabloaguilar/bucketo"
)

func TestBucket_Consume(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		capacity        int64
		tokensToConsume int64
		expectedReturn  bool
	}{
		{
			name:            "should consume tokens when the desired amount is available",
			capacity:        10,
			tokensToConsume: 2,
			expectedReturn:  true,
		},
		{
			name:            "should consume tokens when the desired amount is equal to the available",
			capacity:        10,
			tokensToConsume: 10,
			expectedReturn:  true,
		},
		{
			name:            "should not consume tokens when the desired amount is not available",
			capacity:        10,
			tokensToConsume: 11,
			expectedReturn:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// SETUP
			bucket := bucketo.NewBucket(
				tc.capacity,
				bucketo.WithConsumeStrategy(bucketo.NewDynamicConsume()),
			)

			// TEST
			ok, err := bucket.Consume(tc.tokensToConsume)

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, tc.expectedReturn, ok)
		})
	}
}
