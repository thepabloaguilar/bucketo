package bucketo

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

type Bucket struct {
	capacity      int64
	tokensChannel chan int64

	storage         TokensStorage
	consumeStrategy ConsumeStrategy
	refillers       []Refiller

	m sync.RWMutex
}

func NewBucket(capacity int64, opts ...BucketOpt) *Bucket {
	staticConsume, _ := NewStaticConsume(1) //nolint
	b := &Bucket{
		capacity:        capacity,
		storage:         NewInMemoryTokenStorage(capacity),
		tokensChannel:   make(chan int64, 20),
		consumeStrategy: staticConsume,
		m:               sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Bucket) AvailableTokens() (int64, error) {
	b.m.RLock()
	defer b.m.RUnlock()

	actualTokens, err := b.storage.GetTokens()
	if err != nil {
		return 0, err
	}

	return actualTokens, nil
}

func (b *Bucket) Consume(args any) (bool, error) {
	b.m.Lock()
	defer b.m.Unlock()

	tokensToConsume, err := b.consumeStrategy(args)
	if err != nil {
		return false, err
	}

	actualTokens, err := b.storage.GetTokens()
	if err != nil {
		return false, err
	}

	if tokensToConsume > actualTokens {
		return false, nil
	}

	actualTokens -= tokensToConsume
	err = b.storage.SetTokens(actualTokens)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *Bucket) AddTokens(tokens int64) error {
	b.m.Lock()
	defer b.m.Unlock()

	actualTokens, err := b.storage.GetTokens()
	if err != nil {
		return err
	}

	actualTokens += tokens
	if actualTokens > b.capacity {
		actualTokens = b.capacity
	}

	err = b.storage.SetTokens(actualTokens)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) Start(ctx context.Context) error {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				break
			case tokens := <-b.tokensChannel:
				if err := b.AddTokens(tokens); err != nil { //nolint
					// We're doing nothing here but we should at least log!!
				}
			}
		}
	}(ctx)

	var err error
	for _, refiller := range b.refillers {
		err = multierr.Append(err, refiller.StartRefiller(ctx, b.tokensChannel))
	}

	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) Stop() error {
	var err error
	for _, refiller := range b.refillers {
		err = multierr.Append(err, refiller.StopRefiller())
	}

	if err != nil {
		return err
	}

	return nil
}
