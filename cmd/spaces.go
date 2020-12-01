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
	Arg      string `json:"arg"`
}

type AlfredItems struct {
	Items []AlfredEntry `json:"items"`
}

type Space struct {
	Name      string
	DisplayID uint
	Id        uint
}

// listCmd Args
var (
	OnlyNamed bool
	JSON      bool
	Alfred    bool
)

// switchCmd Args
var (
	SwitchToSpaceNumber int
	SwitchToSpaceName   string
)

func listSpaces() []Space {
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

	return spaces
}

var spacesCmd = &cli.Command{
	Use:   "spaces",
	Short: "Interact with spaces",
}

var listCmd = &cli.Command{
	Use:   "list",
	Short: "List all spaces",
	Run: func(cmd *cli.Command, args []string) {
		var spaces = listSpaces()

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
					Arg:      fmt.Sprint(space.Id),
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

var switchCmd = &cli.Command{
	Use:   "switch",
	Short: "Switch to a space",
	Run: func(cmd *cli.Command, args []string) {
		var spaces = listSpaces()

		// TODO: Improve this using better lookup, but I'm guessing that iterating over
		// ~16 spaces is going to be close enough to O(1) that it doesn't really matter.
		for _, space := range spaces {
			if SwitchToSpaceNumber != -1 && space.Id == uint(SwitchToSpaceNumber) {
				ts2.MoveToSpaceOnDisplay(space.Id, space.DisplayID)
				return
			}

			if SwitchToSpaceName != "" && space.Name == SwitchToSpaceName {
				ts2.MoveToSpaceOnDisplay(space.Id, space.DisplayID)
				return
			}
		}
	},
}

var nameCmd = &cli.Command{
	Use:   "name",
	Short: "Name the current space",
	Args:  cli.ExactArgs(1),
	Run: func(cmd *cli.Command, args []string) {
		// Hardcode to the first display for now
		display := ts2.DisplayList()[0]
		space := ts2.CurrentSpaceNumberOnDisplay(display.ID)
		ts2.SetNameForSpaceOnDisplay(space, args[0], display.ID)
	},
}

func init() {
	spacesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&OnlyNamed, "only-named", "n", false, "Only display named spaces")
	listCmd.Flags().BoolVarP(&JSON, "json", "j", false, "JSON output")
	listCmd.Flags().BoolVarP(&Alfred, "alfred", "a", false, "Output for Alfred's script filter")

	spacesCmd.AddCommand(switchCmd)
	switchCmd.Flags().IntVarP(&SwitchToSpaceNumber, "number", "i", -1, "Switch to a space by number (starts at 1; doesn't currently support multiple displays)")
	switchCmd.Flags().StringVarP(&SwitchToSpaceName, "name", "n", "", "Switch to a space by name")

	spacesCmd.AddCommand(nameCmd)
}
