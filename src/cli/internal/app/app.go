package app

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/ShiunduZachariah/azopscli/internal/azure"
	"github.com/ShiunduZachariah/azopscli/internal/config"
)

type App struct {
	Config         config.Config
	ResourceGroups azure.ResourceGroupLister
}

func New(ctx context.Context) (*App, error) {
	_ = ctx

	if err := config.LoadEnvironmentFiles(); err != nil {
		return nil, fmt.Errorf("load environment files: %w", err)
	}

	cfg := config.Load()
	cfg.Normalize()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("create default azure credential: %w", err)
	}

	return &App{
		Config:         cfg,
		ResourceGroups: azure.NewResourceGroupsClient(credential, nil),
	}, nil
}
