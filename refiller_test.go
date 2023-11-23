package bucketo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/thepabloaguilar/bucketo"
)

func TestTimeRefiller_Success(t *testing.T) {
	t.Parallel()

	t.Run("should send to the tokens channel correctly", func(t *testing.T) {
		// SETUP
		ctx := context.Background()
		tokensChannel := make(chan int64, 10)
		timeToRefill := time.Millisecond * 500
		refiller := bucketo.NewTimeRefiller(10, timeToRefill)

		// TEST
		err := refiller.StartRefiller(ctx, tokensChannel)
		require.NoError(t, err)

		time.Sleep(timeToRefill)

		err = refiller.StopRefiller()
		require.NoError(t, err)

		// ASSERT
		require.Len(t, tokensChannel, 1)
		require.Equal(t, int64(10), <-tokensChannel)
	})

	t.Run("should not send to the tokens channel when we called stop", func(t *testing.T) {
		// SETUP
		ctx := context.Background()
		tokensChannel := make(chan int64, 10)
		timeToRefill := time.Millisecond * 500
		refiller := bucketo.NewTimeRefiller(10, timeToRefill)

		// TEST
		err := refiller.StartRefiller(ctx, tokensChannel)
		require.NoError(t, err)

		// Let send tokens at once
		time.Sleep(timeToRefill)

		err = refiller.StopRefiller()
		require.NoError(t, err)

		// Wait again the time to refill twice, this time it should not send any more tokens
		time.Sleep(timeToRefill * 2)

		// ASSERT
		require.Len(t, tokensChannel, 1)
		require.Equal(t, int64(10), <-tokensChannel)
	})

	t.Run("should not send to the tokens channel when context is canceled", func(t *testing.T) {
		// SETUP
		ctx, cancelCtx := context.WithCancel(context.Background())
		t.Cleanup(cancelCtx)

		tokensChannel := make(chan int64, 10)
		timeToRefill := time.Millisecond * 500
		refiller := bucketo.NewTimeRefiller(10, timeToRefill)

		// TEST
		err := refiller.StartRefiller(ctx, tokensChannel)
		require.NoError(t, err)

		// Let send tokens at once
		time.Sleep(timeToRefill)

		// Canceling context and then waiting for the refill time
		cancelCtx()
		time.Sleep(timeToRefill)

		// ASSERT
		require.Len(t, tokensChannel, 1)
	})
}
