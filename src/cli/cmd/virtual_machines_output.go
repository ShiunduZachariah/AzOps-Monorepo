package cmd

import (
	"fmt"
	"io"

	"github.com/ShiunduZachariah/azopscli/internal/models"
)

func writeVirtualMachines(w io.Writer, vms []models.VirtualMachine, outputFormat string) error {
	switch outputFormat {
	case "", "plain":
		rows := make([][]string, 0, len(vms))
		for _, vm := range vms {
			rows = append(rows, []string{vm.Name, vm.ResourceGroup, vm.Location, vm.VMSize, vm.ID})
		}

		return writeTable(w, []string{"NAME", "RESOURCE GROUP", "LOCATION", "SIZE", "ID"}, rows)
	case "json":
		return writeJSON(w, vms)
	default:
		return fmt.Errorf("unsupported output format %q", outputFormat)
	}
}
