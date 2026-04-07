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

func TestListResourceGroups(t *testing.T) {
	var sawAuth string
	var sawPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		sawPath = r.URL.Path + "?" + r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"value":[{"id":"/subscriptions/sub-123/resourceGroups/rg-one","name":"rg-one","location":"eastus"}]}`))
	}))
	defer server.Close()

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
	if !strings.Contains(sawPath, "/subscriptions/sub-123/resourcegroups") {
		t.Fatalf("unexpected request path: %s", sawPath)
	}
	if len(groups) != 1 || groups[0].Name != "rg-one" || groups[0].Location != "eastus" {
		t.Fatalf("unexpected groups: %#v", groups)
	}
}
