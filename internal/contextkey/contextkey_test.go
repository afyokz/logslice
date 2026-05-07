package contextkey_test

import (
	"context"
	"testing"

	"github.com/yourorg/logslice/internal/contextkey"
)

func TestJobID_StoreAndRetrieve(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextkey.JobID, "abc-123")
	val, ok := ctx.Value(contextkey.JobID).(string)
	if !ok {
		t.Fatal("expected string value for JobID")
	}
	if val != "abc-123" {
		t.Errorf("got %q, want %q", val, "abc-123")
	}
}

func TestInputFile_StoreAndRetrieve(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextkey.InputFile, "/var/log/app.log")
	val, ok := ctx.Value(contextkey.InputFile).(string)
	if !ok {
		t.Fatal("expected string value for InputFile")
	}
	if val != "/var/log/app.log" {
		t.Errorf("got %q, want %q", val, "/var/log/app.log")
	}
}

func TestVerbose_StoreAndRetrieve(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextkey.Verbose, true)
	val, ok := ctx.Value(contextkey.Verbose).(bool)
	if !ok {
		t.Fatal("expected bool value for Verbose")
	}
	if !val {
		t.Error("expected Verbose to be true")
	}
}

func TestKeys_DoNotCollide(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextkey.JobID, "job-1")
	ctx = context.WithValue(ctx, contextkey.InputFile, "file.log")
	ctx = context.WithValue(ctx, contextkey.Verbose, false)

	if v, _ := ctx.Value(contextkey.JobID).(string); v != "job-1" {
		t.Errorf("JobID: got %q, want %q", v, "job-1")
	}
	if v, _ := ctx.Value(contextkey.InputFile).(string); v != "file.log" {
		t.Errorf("InputFile: got %q, want %q", v, "file.log")
	}
	if v, _ := ctx.Value(contextkey.Verbose).(bool); v != false {
		t.Errorf("Verbose: got %v, want false", v)
	}
}

func TestKey_String(t *testing.T) {
	tests := []struct {
		key  fmt.Stringer
		want string
	}{
		{contextkey.JobID, "JobID"},
		{contextkey.InputFile, "InputFile"},
		{contextkey.Verbose, "Verbose"},
	}
	for _, tt := range tests {
		if got := tt.key.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}
