package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newHealthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Verify the CLI is wired correctly",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "ok")
			return err
		},
	}
}
