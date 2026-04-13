package cmd

import (
	"fmt"
	"io"

	"github.com/ShiunduZachariah/azopscli/internal/models"
)

func writeResourceGroups(w io.Writer, groups []models.ResourceGroup, outputFormat string) error {
	switch outputFormat {
	case "", "plain":
		rows := make([][]string, 0, len(groups))
		for _, group := range groups {
			rows = append(rows, []string{group.Name, group.Location, group.ID})
		}

		return writeTable(w, []string{"NAME", "LOCATION", "ID"}, rows)
	case "json":
		return writeJSON(w, groups)
	default:
		return fmt.Errorf("unsupported output format %q", outputFormat)
	}
}
