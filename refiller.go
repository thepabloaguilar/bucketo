package bucketo

import (
	"context"
	"time"
)

// Refiller is likely a background job to refill a bucket automatically.
type Refiller interface {
	StartRefiller(ctx context.Context, tokensChannel chan<- int64) error
	StopRefiller() error
}

type TimeRefiller struct {
	refillRate        int64
	refillRateTime    time.Duration
	lastRefill        time.Time
	internalCtx       context.Context
	internalCtxCancel func()
}

func NewTimeRefiller(refillRate int64, refillRateTime time.Duration) *TimeRefiller {
	return &TimeRefiller{
		refillRate:     refillRate,
		refillRateTime: refillRateTime,
	}
}

// StartRefiller starts a background job and for each time waited iteration it sends the refill
// number to the tokens channel.
func (r *TimeRefiller) StartRefiller(ctx context.Context, tokensChannel chan<- int64) error {
	r.lastRefill = time.Now()

	r.internalCtx, r.internalCtxCancel = context.WithCancel(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(r.refillRateTime)
		for {
			select {
			case <-r.internalCtx.Done():
				ticker.Stop()
				break
			case <-ctx.Done():
				ticker.Stop()
				break
			case <-ticker.C:
				now := time.Now()

				toRefill := now.Sub(r.lastRefill).Milliseconds() / r.refillRateTime.Milliseconds()
				toRefill = toRefill * r.refillRate

				tokensChannel <- toRefill

				r.lastRefill = now
			}
		}
	}(ctx)

	return nil
}

// StopRefiller stops the background job.
// CAUTION: Just call this method once after starting otherwise it'll block your program.
func (r *TimeRefiller) StopRefiller() error {
	if r.internalCtx != nil || r.internalCtx.Err() == nil {
		r.internalCtxCancel()
	}

	return nil
}
