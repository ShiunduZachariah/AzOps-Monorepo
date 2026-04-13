package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/ShiunduZachariah/azopscli/internal/azure"
	"github.com/ShiunduZachariah/azopscli/internal/config"
)

const (
	authModeAuto             = "auto"
	authModeServicePrincipal = "service-principal"
)

type App struct {
	Config          config.Config
	ResourceGroups  azure.ResourceGroupLister
	VirtualMachines azure.VirtualMachineLister
}

func New(ctx context.Context) (*App, error) {
	_ = ctx

	if err := config.LoadEnvironmentFiles(); err != nil {
		return nil, fmt.Errorf("load environment files: %w", err)
	}

	cfg := config.Load()
	cfg.Normalize()

	credential, err := buildCredential(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:          cfg,
		ResourceGroups:  azure.NewResourceGroupsClient(credential, nil),
		VirtualMachines: azure.NewVirtualMachinesClient(credential, nil),
	}, nil
}

func buildCredential(cfg config.Config) (azcore.TokenCredential, error) {
	switch strategy, err := determineCredentialStrategy(cfg); {
	case err != nil:
		return nil, err
	case strategy == authModeServicePrincipal:
		credential, err := azidentity.NewClientSecretCredential(cfg.TenantID, cfg.ClientID, cfg.ClientSecret, nil)
		if err != nil {
			return nil, fmt.Errorf("create service principal credential: %w", err)
		}
		return credential, nil
	default:
		return buildDefaultCredentialChain()
	}
}

func determineCredentialStrategy(cfg config.Config) (string, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.AuthMode)) {
	case "", authModeAuto:
		if cfg.HasServicePrincipalCredentials() {
			return authModeServicePrincipal, nil
		}

		return authModeAuto, nil
	case authModeServicePrincipal:
		if !cfg.HasServicePrincipalCredentials() {
			return "", fmt.Errorf("service principal auth requires tenant id, client id, and client secret")
		}

		return authModeServicePrincipal, nil
	default:
		return "", fmt.Errorf("unsupported auth mode %q: use auto or service-principal", cfg.AuthMode)
	}
}

func buildDefaultCredentialChain() (azcore.TokenCredential, error) {
	credentials := make([]azcore.TokenCredential, 0, 3)

	if credential, err := azidentity.NewAzureCLICredential(nil); err == nil {
		credentials = append(credentials, credential)
	} else {
		return nil, fmt.Errorf("create azure cli credential: %w", err)
	}

	if credential, err := azidentity.NewManagedIdentityCredential(nil); err == nil {
		credentials = append(credentials, credential)
	} else {
		return nil, fmt.Errorf("create managed identity credential: %w", err)
	}

	if credential, err := azidentity.NewAzureDeveloperCLICredential(nil); err == nil {
		credentials = append(credentials, credential)
	} else {
		return nil, fmt.Errorf("create azure developer cli credential: %w", err)
	}

	chain, err := azidentity.NewChainedTokenCredential(credentials, nil)
	if err != nil {
		return nil, fmt.Errorf("create default credential chain: %w", err)
	}

	return chain, nil
}
