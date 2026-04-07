package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ShiunduZachariah/azopscli/internal/azure"
	"github.com/ShiunduZachariah/azopscli/internal/config"
	"github.com/spf13/cobra"
)

type ResourceGroupLister interface {
	ListResourceGroups(ctx context.Context, subscriptionID string) ([]azure.ResourceGroup, error)
}

type Dependencies struct {
	Config         *config.Config
	ResourceGroups ResourceGroupLister
}

func NewRootCommand(deps Dependencies) *cobra.Command {
	cfg := deps.Config
	if cfg == nil {
		defaultCfg := config.Load()
		cfg = &defaultCfg
	}

	cfg.Normalize()

	root := &cobra.Command{
		Use:              "azops",
		Short:            "AzOps CLI for Azure operations",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
	}

	root.PersistentFlags().StringVar(&cfg.SubscriptionID, "subscription-id", cfg.SubscriptionID, "Azure subscription ID (or AZOPS_SUBSCRIPTION_ID / AZURE_SUBSCRIPTION_ID)")
	root.PersistentFlags().StringVar(&cfg.Output, "output", cfg.Output, "Output format: plain or json")

	root.AddCommand(newHealthCommand())
	root.AddCommand(newGroupsCommand(cfg, deps.ResourceGroups))

	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		cfg.Normalize()
		if cfg.Output != "plain" && cfg.Output != "json" {
			return fmt.Errorf("unsupported output format %q: use plain or json", cfg.Output)
		}
		return nil
	}

	return root
}

func newGroupsCommand(cfg *config.Config, lister ResourceGroupLister) *cobra.Command {
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage Azure resource groups",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List Azure resource groups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGroupsList(cmd, cfg, lister)
		},
	}

	groupsCmd.AddCommand(listCmd)
	return groupsCmd
}

func runGroupsList(cmd *cobra.Command, cfg *config.Config, lister ResourceGroupLister) error {
	if lister == nil {
		return fmt.Errorf("resource group lister is not configured")
	}

	subscriptionID := strings.TrimSpace(cfg.SubscriptionID)
	if subscriptionID == "" {
		return fmt.Errorf("subscription id is required: set --subscription-id, AZOPS_SUBSCRIPTION_ID, or AZURE_SUBSCRIPTION_ID")
	}

	groups, err := lister.ListResourceGroups(cmd.Context(), subscriptionID)
	if err != nil {
		return err
	}

	return writeResourceGroups(cmd.OutOrStdout(), groups, cfg.Output)
}
