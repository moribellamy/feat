package flags

import (
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
	ld "github.com/launchdarkly/go-server-sdk/v7"
)

type Factory struct {
	onError func(error)
}

func NewFactory() Factory {
	return Factory{}
}

func (f Factory) OnError(handler func(error)) Factory {
	f.onError = handler
	return f
}

func (f Factory) BoolFlag(key string, defaultValue bool) BoolFlag {
	return BoolFlag{key: key, defaultValue: defaultValue, evalFn: (*ld.LDClient).BoolVariation, onError: f.onError}
}

func (f Factory) StringFlag(key string, defaultValue string) StringFlag {
	return StringFlag{key: key, defaultValue: defaultValue, evalFn: (*ld.LDClient).StringVariation, onError: f.onError}
}

func (f Factory) IntFlag(key string, defaultValue int) IntFlag {
	return IntFlag{key: key, defaultValue: defaultValue, evalFn: (*ld.LDClient).IntVariation, onError: f.onError}
}

func (f Factory) Float64Flag(key string, defaultValue float64) Float64Flag {
	return Float64Flag{key: key, defaultValue: defaultValue, evalFn: (*ld.LDClient).Float64Variation, onError: f.onError}
}

func (f Factory) JSONFlag(key string, defaultValue ldvalue.Value) JSONFlag {
	return JSONFlag{key: key, defaultValue: defaultValue, evalFn: (*ld.LDClient).JSONVariation, onError: f.onError}
}
