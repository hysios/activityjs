// +build !js

package utils

import (
	"io"

	"github.com/Landoop/tableprinter"
)

type Printer struct {
	*tableprinter.Printer
}

func MakeTable(out io.Writer) *Printer {
	printer := tableprinter.New(out)

	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	// printer.HeaderBgColor = tablewriter.BgBlackColor
	// printer.HeaderFgColor = tablewriter.FgGreenColor
	return &Printer{
		Printer: printer,
	}
}
