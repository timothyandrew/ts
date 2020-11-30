package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func print(spaces []Space) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Index", "Display"})

	for _, space := range spaces {
		table.Append([]string{space.Name, fmt.Sprint(space.Id), fmt.Sprint(space.DisplayID)})
	}

	table.Render() // Send output
}
