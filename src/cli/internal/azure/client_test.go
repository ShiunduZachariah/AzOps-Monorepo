package azure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type fakeCredential struct {
	token azcore.AccessToken
	err   error
	calls int
}

func (f *fakeCredential) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	f.calls++
	return f.token, f.err
}

func TestResourceGroupListing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		responseBody string
		wantName     string
		wantLocation string
		wantPath     string
	}{
		{
			name:         "single group",
			responseBody: `{"value":[{"id":"/subscriptions/sub-123/resourceGroups/rg-one","name":"rg-one","location":"eastus"}]}`,
			wantName:     "rg-one",
			wantLocation: "eastus",
			wantPath:     "/subscriptions/sub-123/resourcegroups",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var sawAuth string
			var sawPath string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sawAuth = r.Header.Get("Authorization")
				sawPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			t.Cleanup(server.Close)

			cred := &fakeCredential{
				token: azcore.AccessToken{Token: "token-123", ExpiresOn: time.Now().Add(time.Hour)},
			}
			client := NewResourceGroupsClient(cred, server.Client())
			client.endpoint = server.URL

			groups, err := client.ListResourceGroups(context.Background(), "sub-123")
			if err != nil {
				t.Fatalf("list resource groups: %v", err)
			}

			if cred.calls != 1 {
				t.Fatalf("expected 1 token request, got %d", cred.calls)
			}
			if sawAuth != "Bearer token-123" {
				t.Fatalf("unexpected auth header: %q", sawAuth)
			}
			if !strings.Contains(sawPath, tt.wantPath) {
				t.Fatalf("unexpected request path: %s", sawPath)
			}
			if len(groups) != 1 || groups[0].Name != tt.wantName || groups[0].Location != tt.wantLocation {
				t.Fatalf("unexpected groups: %#v", groups)
			}
		})
	}
}

func TestVirtualMachineListing(t *testing.T) {
	t.Parallel()

	var sawAuth string
	var sawPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		sawPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"value":[{"id":"/subscriptions/sub-123/resourceGroups/rg-one/providers/Microsoft.Compute/virtualMachines/vm-one","name":"vm-one","location":"eastus","properties":{"hardwareProfile":{"vmSize":"Standard_B2s"}}}]}`))
	}))
	t.Cleanup(server.Close)

	cred := &fakeCredential{
		token: azcore.AccessToken{Token: "token-456", ExpiresOn: time.Now().Add(time.Hour)},
	}
	client := NewVirtualMachinesClient(cred, server.Client())
	client.endpoint = server.URL

	vms, err := client.ListVirtualMachines(context.Background(), "sub-123")
	if err != nil {
		t.Fatalf("list virtual machines: %v", err)
	}

	if cred.calls != 1 {
		t.Fatalf("expected 1 token request, got %d", cred.calls)
	}
	if sawAuth != "Bearer token-456" {
		t.Fatalf("unexpected auth header: %q", sawAuth)
	}
	if !strings.Contains(sawPath, "/subscriptions/sub-123/providers/Microsoft.Compute/virtualMachines") {
		t.Fatalf("unexpected request path: %s", sawPath)
	}
	if len(vms) != 1 {
		t.Fatalf("unexpected vms: %#v", vms)
	}
	if vms[0].Name != "vm-one" || vms[0].Location != "eastus" || vms[0].ResourceGroup != "rg-one" || vms[0].VMSize != "Standard_B2s" {
		t.Fatalf("unexpected vm model: %#v", vms[0])
	}
}

func TestListersValidateInputs(t *testing.T) {
	t.Parallel()

	cred := &fakeCredential{token: azcore.AccessToken{Token: "token"}}

	t.Run("resource groups require subscription", func(t *testing.T) {
		client := NewResourceGroupsClient(cred, http.DefaultClient)
		_, err := client.ListResourceGroups(context.Background(), "")
		if err == nil || !strings.Contains(err.Error(), "subscription id is required") {
			t.Fatalf("expected subscription error, got %v", err)
		}
	})

	t.Run("virtual machines require subscription", func(t *testing.T) {
		client := NewVirtualMachinesClient(cred, http.DefaultClient)
		_, err := client.ListVirtualMachines(context.Background(), "")
		if err == nil || !strings.Contains(err.Error(), "subscription id is required") {
			t.Fatalf("expected subscription error, got %v", err)
		}
	})
}
