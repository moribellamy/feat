## Acme

A scratch binary for ad-hoc testing against an existing LaunchDarkly project. Not intended for production use.

```bash
LAUNCHDARKLY_SDK_KEY=your-key go run ./cmd/acme
```

### Required flag keys

The following flags must exist in the LaunchDarkly project associated with the SDK key:

| Key                | Type    | Default        |
|--------------------|---------|----------------|
| `software-version` | float64 | `0.0`          |
| `motd`             | string  | `""`           |
| `beta`             | bool    | `false`        |
| `metadata`         | JSON    | `{}`           |
