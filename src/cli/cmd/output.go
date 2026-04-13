package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func writeTable(w io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, joinRow(headers)); err != nil {
		return err
	}

	for _, row := range rows {
		if _, err := fmt.Fprintln(tw, joinRow(row)); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func joinRow(values []string) string {
	if len(values) == 0 {
		return ""
	}

	row := values[0]
	for _, value := range values[1:] {
		row += "\t" + value
	}

	return row
}
