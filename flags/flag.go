package flags

import (
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
	ld "github.com/launchdarkly/go-server-sdk/v7"
)

type evalFunc[T any] func(client *ld.LDClient, key string, ctx ldcontext.Context, defaultVal T) (T, error)

type Result[T any] struct {
	Value T
	Err   error
}

type Flag[T any] struct {
	key          string
	defaultValue T
	evalFn       evalFunc[T]
	onError      func(error)
}

type BoolFlag = Flag[bool]
type StringFlag = Flag[string]
type IntFlag = Flag[int]
type Float64Flag = Flag[float64]
type JSONFlag = Flag[ldvalue.Value]

func (f Flag[T]) OnError(handler func(error)) Flag[T] {
	f.onError = handler
	return f
}

func (f Flag[T]) Evaluate(client *ld.LDClient, ctx ldcontext.Context) Result[T] {
	value, err := f.evalFn(client, f.key, ctx, f.defaultValue)
	if err != nil && f.onError != nil {
		f.onError(err)
	}
	return Result[T]{Value: value, Err: err}
}
