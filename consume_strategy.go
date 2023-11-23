package bucketo

import (
	"errors"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

var (
	ErrNoExpressionMatched         = errors.New("no expression matched")
	ErrTokensToConsumeNegative     = errors.New("tokens to consume should be positive or zero")
	ErrConsumeArgumentIsNotInteger = errors.New("consume strategy argument is not an integer")
)

// ConsumeStrategy is a function that determines how many tokens should be consumed.
type ConsumeStrategy func(args any) (int64, error)

// NewStaticConsume returns a static consume strategy which means the same token will always be returned.
// For example, NewStaticConsume(1) will always return 1 (one) when called to get the tokens to consume.
func NewStaticConsume(tokens int64) (ConsumeStrategy, error) {
	if tokens < 0 {
		return nil, ErrTokensToConsumeNegative
	}

	return func(_ any) (int64, error) {
		return tokens, nil
	}, nil
}

// NewDynamicConsume returns the same number it receives as argument, if the argument is not an int64
// it'll return an error.
func NewDynamicConsume() ConsumeStrategy {
	return func(args any) (int64, error) {
		argInt, ok := args.(int64)
		if !ok {
			return 0, ErrConsumeArgumentIsNotInteger
		}

		if argInt < 0 {
			return 0, ErrTokensToConsumeNegative
		}

		return argInt, nil
	}
}

type ConsumeExpression struct {
	Expression string
	Tokens     int64

	program *vm.Program
}

// NewExpressionsConsume allows you to dynamically consume tokens based on expressions, the expressions
// can relly on the arguments passed to the strategy, and they always must return a boolean value.
// Expressions are evaluated in the same order from the list, from the first to the last.
//
//	expressions := []ConsumeExpression{
//		{Expression: "my_argument == 1", Tokens: 50},
//		{Expression: `my_second_argument == "a"`, Tokens: 40},
//		{Expression: "true", Tokens: 100}, // This can be considered a default case since it'll always return true
//	}
//
//	consumeStrategy, err := NewExpressionsConsume(expressions)
//	if err != nil {
//		panic(err)
//	}
//
// For more information on how to use/build the expressions see https://expr.medv.io
func NewExpressionsConsume(consumeExpressions []ConsumeExpression) (ConsumeStrategy, error) {
	for idx, expression := range consumeExpressions {
		if expression.Tokens < 0 {
			return nil, ErrTokensToConsumeNegative
		}

		program, err := expr.Compile(
			expression.Expression,
			expr.AsBool(),
			expr.Optimize(true),
		)
		if err != nil {
			return nil, err
		}

		consumeExpressions[idx].program = program
	}

	// VM instance to use for all evaluations, doing this we can have a slight performance increase. Source:
	// https://expr.medv.io/docs/Tips
	v := vm.VM{}

	return func(args any) (int64, error) {
		for _, expression := range consumeExpressions {
			result, err := v.Run(expression.program, args)
			if err != nil {
				return 0, err
			}

			hasMatched := result.(bool) //nolint
			if hasMatched {
				return expression.Tokens, nil
			}
		}

		return 0, ErrNoExpressionMatched
	}, nil
}
