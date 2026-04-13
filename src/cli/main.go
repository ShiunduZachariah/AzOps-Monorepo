package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ShiunduZachariah/azopscli/cmd"
	"github.com/ShiunduZachariah/azopscli/internal/app"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	application, err := app.New(ctx)
	if err != nil {
		return err
	}

	root := cmd.NewRootCommand(cmd.Dependencies{
		Config:          &application.Config,
		ResourceGroups:  application.ResourceGroups,
		VirtualMachines: application.VirtualMachines,
	})
	return root.Execute()
}
