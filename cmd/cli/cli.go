package cli

import (
	"fmt"
	mcli "github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
	"io"
)

type Ui interface {
	mcli.Ui
	mcli.CommandAutocomplete
	Stdout() io.Writer
	Stderr() io.Writer
	ErrorOutput(string)
	HeaderOutput(string)
	WarnOutput(string)
	SuccessOutput(string)
	UnchangedOutput(string)
	Table(tbl *Table)
}

// BasicUI augments mitchellh/cli.BasicUi by exposing the underlying io.Writer.
type BasicUI struct {
	mcli.BasicUi
	mcli.CommandAutocomplete
}

func (b *BasicUI) Stdout() io.Writer {
	return b.BasicUi.Writer
}

func (b *BasicUI) Stderr() io.Writer {
	return b.BasicUi.ErrorWriter
}

func (b *BasicUI) HeaderOutput(s string) {
	b.Output(colorize(fmt.Sprintf("==> %s", s), UiColorNone))
}

func (b *BasicUI) ErrorOutput(s string) {
	b.Output(colorize(fmt.Sprintf(" ! %s", s), UiColorRed))
}

func (b *BasicUI) WarnOutput(s string) {
	b.Output(colorize(fmt.Sprintf(" * %s", s), UiColorYellow))
}

func (b *BasicUI) SuccessOutput(s string) {
	b.Output(colorize(fmt.Sprintf(" ✓ %s", s), UiColorGreen))
}

func (b *BasicUI) UnchangedOutput(s string) {
	b.Output(colorize(fmt.Sprintf("  %s", s), UiColorNone))
}

// Table implements UI.
func (b *BasicUI) Table(tbl *Table) {
	// Build our config and set our options

	table := tablewriter.NewWriter(b.Stdout())

	table.SetHeader(tbl.Headers)
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	for _, row := range tbl.Rows {
		colors := make([]tablewriter.Colors, len(row))
		entries := make([]string, len(row))

		for i, ent := range row {
			entries[i] = ent.Value

			color, ok := colorMapping[ent.Color]
			if ok {
				colors[i] = tablewriter.Colors{color}
			}
		}

		table.Rich(entries, colors)
	}

	table.Render()
}
