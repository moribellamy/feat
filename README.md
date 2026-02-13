# feat

> **feat** /fiːt/ *noun*
> 1. An achievement that requires great courage, skill, or strength.
> 2. *(colloq.)* Short for "feature."

A Go library for type-safe LaunchDarkly feature flag evaluation. Flag keys and defaults are declared once and accessed through typed wrappers, avoiding string duplication and typos.

## Installation

```bash
go get github.com/moribellamy/feat
```

## Usage

```go
import "github.com/moribellamy/feat/flags"

factory := flags.NewFactory().OnError(func(err error) {
    log.Printf("flag error: %v", err)
})

var beta = factory.BoolFlag("beta", false)

result := beta.Evaluate(client, context)
fmt.Println(result.Value) // bool
```

## Error handling

`OnError` can be set at two levels:

1. **Factory** — applies to all flags created by that factory.
2. **Flag** — overrides the factory default for that specific flag.

```go
factory := flags.NewFactory().OnError(func(err error) {
    slog.Warn("flag error", "error", err) // default for all flags
})

var beta = factory.BoolFlag("beta", false)           // uses factory handler
var metadata = factory.JSONFlag("metadata", ldvalue.ObjectBuild().Build()).OnError(func(err error) {
    slog.Error("critical flag failed", "error", err)  // overrides factory handler
})
```

Flags without any error handler (neither factory nor flag-level) will silently return the default value and the error in `Result.Err`.

## Packages

- **`flags`** — Generic `Flag[T]` type, `Factory` for constructing flags with shared defaults, and type aliases (`BoolFlag`, `StringFlag`, `IntFlag`, `Float64Flag`, `JSONFlag`).
- **`cmd/acme`** — A scratch binary for ad-hoc testing against a live LaunchDarkly project. Also serves as a working example of how to integrate the `flags` package. See its [README](cmd/acme/README.md).
