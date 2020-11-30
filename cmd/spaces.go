package cmd

import (
	"encoding/json"
	"fmt"

	cli "github.com/spf13/cobra"
	ts2 "timothyandrew.net/totalspaces/api"
)

type AlfredEntry struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type AlfredItems struct {
	Items []AlfredEntry `json:"items"`
}

type Space struct {
	Name      string
	DisplayID uint
	Id        uint
}

var (
	OnlyNamed bool
	JSON      bool
	Alfred    bool
)

var spacesCmd = &cli.Command{
	Use:   "spaces",
	Short: "Interact with spaces",
}

var listCmd = &cli.Command{
	Use:   "list",
	Short: "List all spaces",
	Run: func(cmd *cli.Command, args []string) {
		var spaces []Space
		var displays = ts2.DisplayList()

		for _, display := range displays {
			var spaceCount = ts2.NumberOfSpacesOnDisplay(display.ID)

			for i := uint(1); i <= spaceCount; i++ {
				var name = ts2.CustomNameForSpaceNumberOnDisplay(i, display.ID)

				if OnlyNamed && name == "" {
					continue
				}

				name = ts2.SpaceNameForSpaceNumberOnDisplay(i, display.ID)

				spaces = append(spaces, Space{
					Name:      name,
					DisplayID: display.ID,
					Id:        i,
				})
			}
		}

		if JSON {
			var output, err = json.Marshal(spaces)
			if err != nil {
				panic("Failed to write JSON")
			}
			fmt.Println(string(output))
		} else if Alfred {
			var entries []AlfredEntry

			for _, space := range spaces {
				entries = append(entries, AlfredEntry{
					Title:    space.Name,
					Subtitle: fmt.Sprintf("Index: %v", space.Id),
				})
			}

			var output, err = json.Marshal(AlfredItems{Items: entries})
			if err != nil {
				panic("Failed to write JSON")
			}

			fmt.Println(string(output))

		} else {
			print(spaces)
		}
	},
}

func init() {
	spacesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&OnlyNamed, "only-named", "n", false, "Only display named spaces")
	listCmd.Flags().BoolVarP(&JSON, "json", "j", false, "JSON output")
	listCmd.Flags().BoolVarP(&Alfred, "alfred", "a", false, "Output for Alfred's script filter")
}
