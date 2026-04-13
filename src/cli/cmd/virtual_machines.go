package cmd

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/ShiunduZachariah/azopscli/internal/azure"
	"github.com/ShiunduZachariah/azopscli/internal/config"
	"github.com/spf13/cobra"
)

func newVirtualMachinesCommand(cfg *config.Config, lister VirtualMachineLister) *cobra.Command {
	vmCmd := &cobra.Command{
		Use:   "vm",
		Short: "Manage Azure virtual machines",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List Azure virtual machines",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVirtualMachinesList(cmd, cfg, lister)
		},
	}

	vmCmd.AddCommand(listCmd)
	return vmCmd
}

func runVirtualMachinesList(cmd *cobra.Command, cfg *config.Config, lister VirtualMachineLister) error {
	if lister == nil {
		var err error
		lister, err = newDefaultVirtualMachineLister()
		if err != nil {
			return err
		}
	}

	subscriptionID := strings.TrimSpace(cfg.SubscriptionID)
	if subscriptionID == "" {
		return fmt.Errorf("subscription id is required: set --subscription-id, AZOPS_SUBSCRIPTION_ID, or AZURE_SUBSCRIPTION_ID")
	}

	vms, err := lister.ListVirtualMachines(cmd.Context(), subscriptionID)
	if err != nil {
		return err
	}

	return writeVirtualMachines(cmd.OutOrStdout(), vms, cfg.Output)
}

func newDefaultVirtualMachineLister() (VirtualMachineLister, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("create default azure credential: %w", err)
	}

	return azure.NewVirtualMachinesClient(credential, nil), nil
}
