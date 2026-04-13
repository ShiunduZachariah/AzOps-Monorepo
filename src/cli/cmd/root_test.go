package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ShiunduZachariah/azopscli/internal/config"
	"github.com/ShiunduZachariah/azopscli/internal/models"
)

type fakeResourceGroupLister struct {
	groups []models.ResourceGroup
	err    error
}

func (f fakeResourceGroupLister) ListResourceGroups(_ context.Context, _ string) ([]models.ResourceGroup, error) {
	return f.groups, f.err
}

type fakeVirtualMachineLister struct {
	vms []models.VirtualMachine
	err error
}

func (f fakeVirtualMachineLister) ListVirtualMachines(_ context.Context, _ string) ([]models.VirtualMachine, error) {
	return f.vms, f.err
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

func TestCLICommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		args            []string
		config          config.Config
		resourceGroups  fakeResourceGroupLister
		virtualMachines fakeVirtualMachineLister
		wantOutput      []string
		wantErr         string
	}{
		{
			name: "groups list json",
			args: []string{"groups", "list"},
			config: config.Config{
				SubscriptionID: "sub-123",
				Output:         "json",
			},
			resourceGroups: fakeResourceGroupLister{
				groups: []models.ResourceGroup{
					{Name: "rg-one", Location: "eastus", ID: "/subscriptions/sub-123/resourceGroups/rg-one"},
				},
			},
			wantOutput: []string{`"name": "rg-one"`, `"location": "eastus"`},
		},
		{
			name: "groups list plain",
			args: []string{"groups", "list"},
			config: config.Config{
				SubscriptionID: "sub-123",
				Output:         "plain",
			},
			resourceGroups: fakeResourceGroupLister{
				groups: []models.ResourceGroup{
					{Name: "rg-one", Location: "eastus", ID: "/subscriptions/sub-123/resourceGroups/rg-one"},
				},
			},
			wantOutput: []string{"NAME", "rg-one", "eastus"},
		},
		{
			name: "vm list json",
			args: []string{"vm", "list"},
			config: config.Config{
				SubscriptionID: "sub-123",
				Output:         "json",
			},
			virtualMachines: fakeVirtualMachineLister{
				vms: []models.VirtualMachine{
					{Name: "vm-one", Location: "eastus", ResourceGroup: "rg-one", VMSize: "Standard_B2s", ID: "/subscriptions/sub-123/resourceGroups/rg-one/providers/Microsoft.Compute/virtualMachines/vm-one"},
				},
			},
			wantOutput: []string{`"name": "vm-one"`, `"resourceGroup": "rg-one"`, `"vmSize": "Standard_B2s"`},
		},
		{
			name: "vm list plain",
			args: []string{"vm", "list"},
			config: config.Config{
				SubscriptionID: "sub-123",
				Output:         "plain",
			},
			virtualMachines: fakeVirtualMachineLister{
				vms: []models.VirtualMachine{
					{Name: "vm-one", Location: "eastus", ResourceGroup: "rg-one", VMSize: "Standard_B2s", ID: "/subscriptions/sub-123/resourceGroups/rg-one/providers/Microsoft.Compute/virtualMachines/vm-one"},
				},
			},
			wantOutput: []string{"RESOURCE GROUP", "vm-one", "Standard_B2s"},
		},
		{
			name: "missing subscription errors for groups",
			args: []string{"groups", "list"},
			config: config.Config{
				Output: "json",
			},
			resourceGroups: fakeResourceGroupLister{},
			wantErr:        "subscription id is required",
		},
		{
			name: "missing subscription errors for vm",
			args: []string{"vm", "list"},
			config: config.Config{
				Output: "json",
			},
			virtualMachines: fakeVirtualMachineLister{},
			wantErr:         "subscription id is required",
		},
		{
			name: "invalid output errors",
			args: []string{"groups", "list"},
			config: config.Config{
				SubscriptionID: "sub-123",
				Output:         "yaml",
			},
			wantErr: "unsupported output format",
		},
		{
			name: "vm lister error bubbles up",
			args: []string{"vm", "list"},
			config: config.Config{
				SubscriptionID: "sub-123",
				Output:         "json",
			},
			virtualMachines: fakeVirtualMachineLister{
				err: errors.New("boom"),
			},
			wantErr: "boom",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := NewRootCommand(Dependencies{
				Config:          &tt.config,
				ResourceGroups:  tt.resourceGroups,
				VirtualMachines: tt.virtualMachines,
			})

			buf := &bytes.Buffer{}
			root.SetOut(buf)
			root.SetErr(buf)
			root.SetArgs(tt.args)

			err := root.Execute()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("execute command: %v", err)
			}

			got := buf.String()
			for _, want := range tt.wantOutput {
				if !strings.Contains(got, want) {
					t.Fatalf("expected output to contain %q, got %s", want, got)
				}
			}
		})
	}
}
