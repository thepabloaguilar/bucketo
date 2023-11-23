package bucketo

type BucketOpt func(b *Bucket)

func WithRefiller(refiller Refiller) BucketOpt {
	return func(b *Bucket) {
		b.refillers = append(b.refillers, refiller)
	}
}

func WithStorage(storage TokensStorage) BucketOpt {
	return func(b *Bucket) {
		b.storage = storage
	}
}

func WithConsumeStrategy(strategy ConsumeStrategy) BucketOpt {
	return func(b *Bucket) {
		b.consumeStrategy = strategy
	}
}
