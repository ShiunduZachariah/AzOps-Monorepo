package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/ShiunduZachariah/azopscli/internal/azure"
)

func writeResourceGroups(w io.Writer, groups []azure.ResourceGroup, outputFormat string) error {
	switch outputFormat {
	case "", "plain":
		tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(tw, "NAME\tLOCATION\tID"); err != nil {
			return err
		}
		for _, group := range groups {
			if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", group.Name, group.Location, group.ID); err != nil {
				return err
			}
		}
		return tw.Flush()
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(groups)
	default:
		return fmt.Errorf("unsupported output format %q", outputFormat)
	}
}
