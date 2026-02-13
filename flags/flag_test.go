package flags

import (
	"errors"
	"testing"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
	ld "github.com/launchdarkly/go-server-sdk/v7"
	"github.com/launchdarkly/go-server-sdk/v7/testhelpers/ldtestdata"
)

type testHarness struct {
	client *ld.LDClient
	td     *ldtestdata.TestDataSource
	ctx    ldcontext.Context
}

func setup(t *testing.T) *testHarness {
	t.Helper()
	td := ldtestdata.DataSource()
	config := ld.Config{
		DataSource: td,
	}
	client, err := ld.MakeCustomClient("fake-key", config, 0)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return &testHarness{
		client: client,
		td:     td,
		ctx:    ldcontext.New("test-user"),
	}
}

func TestBoolFlag_Evaluate(t *testing.T) {
	h := setup(t)
	f := NewFactory()

	t.Run("returns true when flag is true", func(t *testing.T) {
		h.td.Update(h.td.Flag("my-bool").BooleanFlag().VariationForAll(true))
		result := f.BoolFlag("my-bool", false).Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if result.Value != true {
			t.Errorf("got %v, want true", result.Value)
		}
	})

	t.Run("returns false when flag is false", func(t *testing.T) {
		h.td.Update(h.td.Flag("my-bool").BooleanFlag().VariationForAll(false))
		result := f.BoolFlag("my-bool", true).Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if result.Value != false {
			t.Errorf("got %v, want false", result.Value)
		}
	})

	t.Run("returns default when flag is unknown", func(t *testing.T) {
		result := f.BoolFlag("nonexistent", true).Evaluate(h.client, h.ctx)
		if result.Value != true {
			t.Errorf("got %v, want true (default)", result.Value)
		}
	})
}

func TestStringFlag_Evaluate(t *testing.T) {
	h := setup(t)
	f := NewFactory()

	t.Run("returns configured value", func(t *testing.T) {
		h.td.Update(h.td.Flag("my-string").ValueForAll(ldvalue.String("hello")))
		result := f.StringFlag("my-string", "default").Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if result.Value != "hello" {
			t.Errorf("got %q, want %q", result.Value, "hello")
		}
	})

	t.Run("returns default when flag is unknown", func(t *testing.T) {
		result := f.StringFlag("nonexistent-string", "fallback").Evaluate(h.client, h.ctx)
		if result.Value != "fallback" {
			t.Errorf("got %q, want %q", result.Value, "fallback")
		}
	})
}

func TestIntFlag_Evaluate(t *testing.T) {
	h := setup(t)
	f := NewFactory()

	t.Run("returns configured value", func(t *testing.T) {
		h.td.Update(h.td.Flag("my-int").ValueForAll(ldvalue.Int(42)))
		result := f.IntFlag("my-int", 0).Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if result.Value != 42 {
			t.Errorf("got %d, want %d", result.Value, 42)
		}
	})

	t.Run("returns default when flag is unknown", func(t *testing.T) {
		result := f.IntFlag("nonexistent-int", 99).Evaluate(h.client, h.ctx)
		if result.Value != 99 {
			t.Errorf("got %d, want %d", result.Value, 99)
		}
	})
}

func TestFloat64Flag_Evaluate(t *testing.T) {
	h := setup(t)
	f := NewFactory()

	t.Run("returns configured value", func(t *testing.T) {
		h.td.Update(h.td.Flag("my-float").ValueForAll(ldvalue.Float64(3.14)))
		result := f.Float64Flag("my-float", 0.0).Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if result.Value != 3.14 {
			t.Errorf("got %f, want %f", result.Value, 3.14)
		}
	})

	t.Run("returns default when flag is unknown", func(t *testing.T) {
		result := f.Float64Flag("nonexistent-float", 2.72).Evaluate(h.client, h.ctx)
		if result.Value != 2.72 {
			t.Errorf("got %f, want %f", result.Value, 2.72)
		}
	})
}

func TestJSONFlag_Evaluate(t *testing.T) {
	h := setup(t)
	f := NewFactory()

	t.Run("returns configured value", func(t *testing.T) {
		expected := ldvalue.ObjectBuild().Set("key", ldvalue.String("val")).Build()
		h.td.Update(h.td.Flag("my-json").ValueForAll(expected))
		result := f.JSONFlag("my-json", ldvalue.Null()).Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if !result.Value.Equal(expected) {
			t.Errorf("got %s, want %s", result.Value.JSONString(), expected.JSONString())
		}
	})

	t.Run("returns default when flag is unknown", func(t *testing.T) {
		dflt := ldvalue.String("fallback")
		result := f.JSONFlag("nonexistent-json", dflt).Evaluate(h.client, h.ctx)
		if !result.Value.Equal(dflt) {
			t.Errorf("got %s, want %s", result.Value.JSONString(), dflt.JSONString())
		}
	})
}

func TestFlag_OnError(t *testing.T) {
	h := setup(t)

	t.Run("called on evaluation error", func(t *testing.T) {
		var captured error
		flag := NewFactory().BoolFlag("nonexistent", false).OnError(func(err error) {
			captured = err
		})
		result := flag.Evaluate(h.client, h.ctx)
		if result.Err == nil {
			t.Fatal("expected an error for unknown flag")
		}
		if captured == nil {
			t.Fatal("onError handler was not called")
		}
		if !errors.Is(captured, result.Err) {
			t.Errorf("onError received %v, want %v", captured, result.Err)
		}
	})

	t.Run("not called on success", func(t *testing.T) {
		h.td.Update(h.td.Flag("exists").BooleanFlag().VariationForAll(true))
		called := false
		flag := NewFactory().BoolFlag("exists", false).OnError(func(err error) {
			called = true
		})
		result := flag.Evaluate(h.client, h.ctx)
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if called {
			t.Error("onError handler should not be called on success")
		}
	})
}

func TestFactory_OnError_PropagatedToFlags(t *testing.T) {
	h := setup(t)

	var captured error
	f := NewFactory().OnError(func(err error) {
		captured = err
	})

	flag := f.BoolFlag("nonexistent", false)
	result := flag.Evaluate(h.client, h.ctx)
	if result.Err == nil {
		t.Fatal("expected an error for unknown flag")
	}
	if captured == nil {
		t.Fatal("factory-level onError was not propagated to flag")
	}
	if !errors.Is(captured, result.Err) {
		t.Errorf("onError received %v, want %v", captured, result.Err)
	}
}

func TestFlag_OnError_Override(t *testing.T) {
	h := setup(t)

	factoryCalled := false
	f := NewFactory().OnError(func(err error) {
		factoryCalled = true
	})

	var captured error
	flag := f.BoolFlag("nonexistent", false).OnError(func(err error) {
		captured = err
	})

	result := flag.Evaluate(h.client, h.ctx)
	if result.Err == nil {
		t.Fatal("expected an error for unknown flag")
	}
	if factoryCalled {
		t.Error("factory-level onError should not be called when overridden")
	}
	if captured == nil {
		t.Fatal("flag-level onError was not called")
	}
}

func TestFlag_DefaultValue_ReturnedWhenKeyMissing(t *testing.T) {
	h := setup(t)
	f := NewFactory()

	t.Run("bool", func(t *testing.T) {
		result := f.BoolFlag("missing-bool", true).Evaluate(h.client, h.ctx)
		if result.Value != true {
			t.Errorf("got %v, want true", result.Value)
		}
	})

	t.Run("string", func(t *testing.T) {
		result := f.StringFlag("missing-string", "def").Evaluate(h.client, h.ctx)
		if result.Value != "def" {
			t.Errorf("got %q, want %q", result.Value, "def")
		}
	})

	t.Run("int", func(t *testing.T) {
		result := f.IntFlag("missing-int", 7).Evaluate(h.client, h.ctx)
		if result.Value != 7 {
			t.Errorf("got %d, want %d", result.Value, 7)
		}
	})

	t.Run("float64", func(t *testing.T) {
		result := f.Float64Flag("missing-float", 1.5).Evaluate(h.client, h.ctx)
		if result.Value != 1.5 {
			t.Errorf("got %f, want %f", result.Value, 1.5)
		}
	})

	t.Run("json", func(t *testing.T) {
		dflt := ldvalue.String("default-json")
		result := f.JSONFlag("missing-json", dflt).Evaluate(h.client, h.ctx)
		if !result.Value.Equal(dflt) {
			t.Errorf("got %s, want %s", result.Value.JSONString(), dflt.JSONString())
		}
	})
}
