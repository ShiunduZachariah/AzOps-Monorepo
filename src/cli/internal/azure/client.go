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
	"github.com/ShiunduZachariah/azopscli/internal/models"
)

const defaultManagementEndpoint = "https://management.azure.com"

const (
	resourceGroupsAPIVersion  = "2021-04-01"
	virtualMachinesAPIVersion = "2025-04-01"
)

type ResourceGroup = models.ResourceGroup
type VirtualMachine = models.VirtualMachine

type ResourceGroupLister interface {
	ListResourceGroups(ctx context.Context, subscriptionID string) ([]ResourceGroup, error)
}

type VirtualMachineLister interface {
	ListVirtualMachines(ctx context.Context, subscriptionID string) ([]VirtualMachine, error)
}

type ResourceGroupsClient struct {
	credential azcore.TokenCredential
	httpClient *http.Client
	endpoint   string
}

type VirtualMachinesClient struct {
	credential azcore.TokenCredential
	httpClient *http.Client
	endpoint   string
}

func NewResourceGroupsClient(credential azcore.TokenCredential, httpClient *http.Client) *ResourceGroupsClient {
	return &ResourceGroupsClient{
		credential: credential,
		httpClient: newHTTPClient(httpClient),
		endpoint:   defaultManagementEndpoint,
	}
}

func NewVirtualMachinesClient(credential azcore.TokenCredential, httpClient *http.Client) *VirtualMachinesClient {
	return &VirtualMachinesClient{
		credential: credential,
		httpClient: newHTTPClient(httpClient),
		endpoint:   defaultManagementEndpoint,
	}
}

func newHTTPClient(httpClient *http.Client) *http.Client {
	if httpClient != nil {
		return httpClient
	}

	return &http.Client{Timeout: 30 * time.Second}
}

func (c *ResourceGroupsClient) ListResourceGroups(ctx context.Context, subscriptionID string) ([]ResourceGroup, error) {
	if c == nil {
		return nil, fmt.Errorf("resource groups client is nil")
	}

	var payload struct {
		Value []ResourceGroup `json:"value"`
	}

	if err := c.doList(ctx, subscriptionID, fmt.Sprintf("/subscriptions/%s/resourcegroups", url.PathEscape(strings.TrimSpace(subscriptionID))), resourceGroupsAPIVersion, &payload); err != nil {
		return nil, err
	}

	return payload.Value, nil
}

func (c *VirtualMachinesClient) ListVirtualMachines(ctx context.Context, subscriptionID string) ([]VirtualMachine, error) {
	if c == nil {
		return nil, fmt.Errorf("virtual machines client is nil")
	}

	var payload struct {
		Value []virtualMachineResponse `json:"value"`
	}

	if err := c.doList(ctx, subscriptionID, fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Compute/virtualMachines", url.PathEscape(strings.TrimSpace(subscriptionID))), virtualMachinesAPIVersion, &payload); err != nil {
		return nil, err
	}

	vms := make([]VirtualMachine, 0, len(payload.Value))
	for _, item := range payload.Value {
		vms = append(vms, item.toModel())
	}

	return vms, nil
}

func (c *ResourceGroupsClient) doList(ctx context.Context, subscriptionID, resourcePath, apiVersion string, target any) error {
	return doList(ctx, c.credential, c.httpClient, c.endpoint, subscriptionID, resourcePath, apiVersion, target)
}

func (c *VirtualMachinesClient) doList(ctx context.Context, subscriptionID, resourcePath, apiVersion string, target any) error {
	return doList(ctx, c.credential, c.httpClient, c.endpoint, subscriptionID, resourcePath, apiVersion, target)
}

func doList(ctx context.Context, credential azcore.TokenCredential, httpClient *http.Client, endpoint, subscriptionID, resourcePath, apiVersion string, target any) error {
	if credential == nil {
		return fmt.Errorf("credential is required")
	}

	subscriptionID = strings.TrimSpace(subscriptionID)
	if subscriptionID == "" {
		return fmt.Errorf("subscription id is required")
	}

	token, err := credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return fmt.Errorf("acquire azure token: %w", err)
	}

	requestURL := fmt.Sprintf("%s%s?api-version=%s", strings.TrimRight(endpoint, "/"), resourcePath, apiVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("build azure request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call azure resource manager: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("azure resource manager returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode azure response: %w", err)
	}

	return nil
}

type virtualMachineResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Properties struct {
		HardwareProfile struct {
			VMSize string `json:"vmSize"`
		} `json:"hardwareProfile"`
	} `json:"properties"`
}

func (r virtualMachineResponse) toModel() VirtualMachine {
	return VirtualMachine{
		ID:            r.ID,
		Name:          r.Name,
		Location:      r.Location,
		ResourceGroup: resourceGroupFromID(r.ID),
		VMSize:        r.Properties.HardwareProfile.VMSize,
	}
}

func resourceGroupFromID(id string) string {
	parts := strings.Split(id, "/")
	for i := 0; i < len(parts)-1; i++ {
		if strings.EqualFold(parts[i], "resourceGroups") {
			return parts[i+1]
		}
	}

	return ""
}
