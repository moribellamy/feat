package main

import (
	"log/slog"

	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"

	"github.com/moribellamy/feat/flags"
)

type FeatureFlags struct {
	SoftwareVersion flags.Float64Flag
	Motd            flags.StringFlag
	Beta            flags.BoolFlag
	Metadata        flags.JSONFlag
}

var featureFlags = func() FeatureFlags {
	factory := flags.NewFactory().OnError(func(err error) {
		slog.Warn("flag evaluation failed", "error", err)
	})
	return FeatureFlags{
		SoftwareVersion: factory.Float64Flag("software-version", 0.0),
		Motd:            factory.StringFlag("motd", ""),
		Beta:            factory.BoolFlag("beta", false),
		Metadata: factory.JSONFlag("metadata", ldvalue.ObjectBuild().Build()).OnError(func(err error) {
			slog.Error("metadata evaluation is really important, and it failed", "error", err)
		}),
	}
}()
