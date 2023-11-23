# bucketo

[![test](https://github.com/thepabloaguilar/bucketo/actions/workflows/test.yaml/badge.svg)](https://github.com/thepabloaguilar/bucketo/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/thepabloaguilar/bucketo/branch/main/graph/badge.svg?token=qFlORZnn09)](https://codecov.io/gh/thepabloaguilar/bucketo)
[![Go Report Card](https://goreportcard.com/badge/github.com/thepabloaguilar/bucketo)](https://goreportcard.com/report/github.com/thepabloaguilar/bucketo)

## Introduction

bucketo is a package to help you building [token buckets](https://en.wikipedia.org/wiki/Token_bucket), it's designed to be
very flexible yet simple! You can just start using the following code:

```go
package main

import (
	"fmt"
	"github.com/thepabloaguilar/bucketo"
)

func main() {
	// Creates a bucket with 10 of capacity
	bucket := bucketo.NewBucket(10)

	// Consumes 1 token from the bucket, we're ignoring the error return, but you shouldn't
	ok, _ := bucket.Consume(1)
	fmt.Println("Was consumed ->", ok) // true
	
	currentTokens, _ := bucket.AvailableTokens()
	fmt.Println("Current tokens ->", currentTokens) // 9
	
	// Try to consume 10 tokens, but we have only 9 available
	ok, _ = bucket.Consume(10)
	fmt.Println("Was consumed ->", ok) // false

	currentTokens, _ = bucket.AvailableTokens()
	fmt.Println("Current tokens ->", currentTokens) // 9
}
```

For more information on how to configure you're bucket read the next sessions.

## Features

### Tokens Storage

Tokens Storage is the interface where your bucket get and set the tokens amount, so far this package only
provides one token storage but you can implement other storages by respecting the `TokenStorage` interface.

#### InMemoryTokenStorage

This is the default token storage when you don't specify any, it just stores the token amount in memory.

```go
package main

import (
	"fmt"
	"github.com/thepabloaguilar/bucketo"
)

func main() {
	bucketCapacity := int64(100)
	bucket := bucketo.NewBucket(
		bucketCapacity,
		bucketo.WithStorage(bucketo.NewInMemoryTokenStorage(bucketCapacity)), // You can also pass it explicit
	)

	fmt.Println(bucket.AvailableTokens()) // 100
}
```

### Consume Strategy

All the meaning of having a bucket is to consume it, the way you can consume the tokens change depending on your requirements.
By default, the bucket uses the [Static Consume](#consume-strategy) strategy where the static consume number is 1.

#### Static Consume

Static consume strategy will return the same number every time the _Consume_ method is called.

```go
package main

import (
	"fmt"
	"github.com/thepabloaguilar/bucketo"
)

func main() {
	consumeStrategy, _ := bucketo.NewStaticConsume(2) // It will always return 2 as the token number to consume
	bucket := bucketo.NewBucket(100, bucketo.WithConsumeStrategy(consumeStrategy))

	fmt.Println(bucket.AvailableTokens()) // 100

	// Note we're passing `nil` to the Consume method because the consume strategy does not need any info.
	_, _ = bucket.Consume(nil)
	fmt.Println(bucket.AvailableTokens()) // 98
}
```

#### Dynamic Consume

Sometimes we want to consume a non-static number of tokens depending on some variables,
this is where the Dynamic Consume comes to play as it allows we to pass the number of tokens to consume.

```go
package main

import (
	"fmt"
	"github.com/thepabloaguilar/bucketo"
)

func main() {
	bucket := bucketo.NewBucket(
		100,
		bucketo.WithConsumeStrategy(bucketo.NewDynamicConsume()),
	)

	fmt.Println(bucket.AvailableTokens()) // 100

	// Note we're passing `10` to the Consume method, this number is passed to the consume strategy.
	_, _ = bucket.Consume(10)
	fmt.Println(bucket.AvailableTokens()) // 90
}
```

#### Expression Consume

When you want to dynamic consume tokens and have expressions to express your requirements you can delegate
this part to the Expression Consume strategy implementation. Let's see how it works:

```go
package main

import (
	"fmt"
	"github.com/thepabloaguilar/bucketo"
)

func main() {
	// The expressions it'll be used by the consume strategy
	expressions := []bucketo.ConsumeExpression{
		{
			Expression: "my_argument == 1",
			Tokens:     50,
		},
		{
			Expression: `my_second_argument == "a"`,
			Tokens:     40,
		},
		{
			Expression: "true", // This can be considered a default case since it'll always return true
			Tokens:     10,
		},
	}
	consumeStrategy, _ := bucketo.NewExpressionsConsume(expressions)
	bucket := bucketo.NewBucket(100, bucketo.WithConsumeStrategy(consumeStrategy))

	fmt.Println(bucket.AvailableTokens()) // 100

	// Note we're passing a map to the Consume method which is passed to the consume strategy.
	// The value we pass here will be used as variables to evaluate the expressions, here we're
	// passing a map but it could be a struct.
	// The consume strategy will match the first expression in the list.
	//
	// For more information on how to use/build the expressions see https://expr.medv.io
	_, _ = bucket.Consume(map[string]int{
		"my_argument": 1,
	})
	fmt.Println(bucket.AvailableTokens()) // 50

	// Passing an empty map here will make all the expressions failing but the last, our default one.
	_, _ = bucket.Consume(map[string]int{})
	fmt.Println(bucket.AvailableTokens()) // 40

	// It'll match the second expression
	_, _ = bucket.Consume(map[string]string{
		"my_second_argument": "a",
	})
	fmt.Println(bucket.AvailableTokens()) // 0
}
```

### Refillers

The other most important part of a bucket is also to return the tokens to our bucket, you can do it manually
by calling the `AddTokens` method but you can also automate it. This package provides a time refiller.

#### Time Refiller

The Time Refiller will notify the bucket to add more tokens periodically, we can specify the time rate we will
refill and how much to refill.

```go
package main

import (
	"context"
	"fmt"
	"github.com/thepabloaguilar/bucketo"
	"time"
)

func main() {
	consumeStrategy, _ := bucketo.NewStaticConsume(10)
	timeRefiller := bucketo.NewTimeRefiller(10, time.Second) // It'll refill 10 tokens per second
	bucket := bucketo.NewBucket(
		100,
		bucketo.WithConsumeStrategy(consumeStrategy),
		bucketo.WithRefiller(timeRefiller),
	)

	// IMPORTANT: IT'S VERY IMPORTANT TO CALL THE START METHOD
	// The Start method start all the refillers and start listening to them.
	_ = bucket.Start(context.Background())

	fmt.Println(bucket.AvailableTokens()) // 90

	// After 1 second
	fmt.Println(bucket.AvailableTokens()) // 100
}
```
