package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/ShiunduZachariah/azopscli/internal/azure"
	"github.com/ShiunduZachariah/azopscli/internal/config"
)

type fakeLister struct {
	groups []azure.ResourceGroup
	err    error
}

func (f fakeLister) ListResourceGroups(_ context.Context, _ string) ([]azure.ResourceGroup, error) {
	return f.groups, f.err
}

func TestHealthCommand(t *testing.T) {
	root := NewRootCommand(Dependencies{Config: &config.Config{}})
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"health"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute health: %v", err)
	}

	if got := strings.TrimSpace(buf.String()); got != "ok" {
		t.Fatalf("unexpected health output: %q", got)
	}
}

func TestGroupsListJSON(t *testing.T) {
	root := NewRootCommand(Dependencies{
		Config: &config.Config{
			SubscriptionID: "sub-123",
			Output:         "json",
		},
		ResourceGroups: fakeLister{
			groups: []azure.ResourceGroup{
				{Name: "rg-one", Location: "eastus", ID: "/subscriptions/sub-123/resourceGroups/rg-one"},
			},
		},
	})

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"groups", "list"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute groups list: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, `"name": "rg-one"`) || !strings.Contains(got, `"location": "eastus"`) {
		t.Fatalf("unexpected json output: %s", got)
	}
}
