// +build !js

package context

import (
	"fmt"
	"io"

	// "github.com/WeiZhang555/tabwriter"

	"github.com/juju/ansiterm/tabwriter"

	"github.com/olekukonko/tablewriter"
)

func Table(out io.Writer) func(ctx Context) {
	return func(ctx Context) {
		const padding = 3

		item := ctx.Item()
		w := tabwriter.NewWriter(out, 0, 0, padding, '-', tabwriter.AlignRight|tabwriter.Debug)

		if item != nil {
			fmt.Fprintf(w, "%s\t%s\t%0.2f\t%d\t%0.2f\n", item.ID, item.Title, item.Price, item.Quantity, item.Subtotal)
		}
		w.Flush()
		fmt.Fprintln(w)
		order := ctx.Order()
		w = tabwriter.NewWriter(out, 0, 0, padding, '-', tabwriter.AlignRight|tabwriter.Debug)

		if order != nil {
			for _, item := range order.Items {
				fmt.Fprintf(w, "%s\t%s\t%0.2f\t%d\t%0.2f\n", item.ID, item.Title, item.Price, item.Quantity, item.Subtotal)
			}
		}
		w.Flush()
	}
}

func TableEx(out io.Writer) func(ctx Context) {
	return func(ctx Context) {
		item := ctx.Item()

		if item != nil {
			table := tablewriter.NewWriter(out)
			table.SetHeader([]string{"ID", "Title", "Price", "Quantity", "Subtotal"})

			table.Append([]string{item.ID.String(), item.Title, fmt.Sprintf("%0.2f", item.Price), fmt.Sprintf("%d", item.Quantity), fmt.Sprintf("%0.2f", item.Subtotal)})
			table.Render()
		}
		order := ctx.Order()

		if order != nil {
			table := tablewriter.NewWriter(out)
			table.SetHeader([]string{"ID", "Title", "Price", "Quantity", "Subtotal"})
			for _, item := range order.Items {
				table.Append([]string{item.ID.String(), item.Title, fmt.Sprintf("%0.2f", item.Price), fmt.Sprintf("%d", item.Quantity), fmt.Sprintf("%0.2f", item.Subtotal)})
			}
			table.SetFooter([]string{"", "", "Total", fmt.Sprintf("% 5d", order.Count), fmt.Sprintf("%0.2f", order.Total)}) // Add Footer

			table.Render()

		}
	}
}
