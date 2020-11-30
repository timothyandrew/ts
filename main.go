package main

import (
	"fmt"

	totalspaces2 "timothyandrew.net/totalspaces/api"
)

func main() {
	var displays = totalspaces2.DisplayList()

	for _, display := range displays {
		fmt.Printf("%v\n", display)
	}
}
