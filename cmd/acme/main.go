package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v7"
)

func main() {
	sdkKey := os.Getenv("LAUNCHDARKLY_SDK_KEY")
	if sdkKey == "" {
		slog.Error("LAUNCHDARKLY_SDK_KEY environment variable is required")
		os.Exit(1)
	}

	client, err := ld.MakeClient(sdkKey, 5*time.Second)
	if err != nil {
		slog.Error("failed to create LaunchDarkly client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	context := ldcontext.New("example-user")

	svResult := featureFlags.SoftwareVersion.Evaluate(client, context)
	slog.Info("evaluated flag", "software-version", svResult.Value)
	if svResult.Err != nil {
		slog.Error("failed to evaluate software version", "error", svResult.Err)
	}

	motdResult := featureFlags.Motd.Evaluate(client, context)
	slog.Info("evaluated flag", "motd", motdResult.Value)

	betaResult := featureFlags.Beta.Evaluate(client, context)
	slog.Info("evaluated flag", "beta", betaResult.Value)

	metadataResult := featureFlags.Metadata.Evaluate(client, context)
	slog.Info("evaluated flag", "metadata", metadataResult.Value.JSONString())
}
