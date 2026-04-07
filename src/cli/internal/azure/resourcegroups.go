package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const defaultManagementEndpoint = "https://management.azure.com"
const resourceGroupsAPIVersion = "2021-04-01"

type ResourceGroup struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags,omitempty"`
}

type ResourceGroupsClient struct {
	credential azcore.TokenCredential
	httpClient *http.Client
	endpoint   string
}

type ResourceGroupLister interface {
	ListResourceGroups(ctx context.Context, subscriptionID string) ([]ResourceGroup, error)
}

func NewResourceGroupsClient(credential azcore.TokenCredential, httpClient *http.Client) *ResourceGroupsClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &ResourceGroupsClient{
		credential: credential,
		httpClient: httpClient,
		endpoint:   defaultManagementEndpoint,
	}
}

func (c *ResourceGroupsClient) ListResourceGroups(ctx context.Context, subscriptionID string) ([]ResourceGroup, error) {
	if c == nil {
		return nil, fmt.Errorf("resource groups client is nil")
	}
	if c.credential == nil {
		return nil, fmt.Errorf("credential is required")
	}
	subscriptionID = strings.TrimSpace(subscriptionID)
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription id is required")
	}

	token, err := c.credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("acquire azure token: %w", err)
	}

	endpoint := strings.TrimRight(c.endpoint, "/")
	resourceGroupsURL := fmt.Sprintf("%s/subscriptions/%s/resourcegroups?api-version=%s", endpoint, url.PathEscape(subscriptionID), resourceGroupsAPIVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, resourceGroupsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build resource groups request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call azure resource groups api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("azure resource groups api returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		Value []ResourceGroup `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode resource groups response: %w", err)
	}

	return payload.Value, nil
}
