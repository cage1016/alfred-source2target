/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2go/alfred"
)

// t2fCmd represents the t2f command
var t2fCmd = &cobra.Command{
	Use:   "t2f",
	Short: "Targets 2 find",
	Run:   runT2fCmd,
}

func runT2fCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingTargets(wf)

	for name, path := range data {
		wi := wf.NewItem(name).
			Subtitle(fmt.Sprintf("⌘ ,↩ Move / Copy files from '%s' to frontmost Finder", path)).
			Valid(true).
			Arg(path)

		wi.Cmd().
			Subtitle("↩ Enter Action menu to Add / Remove target folder").
			Valid(true)
	}

	if len(data) == 0 {
		wf.NewItem("ADD").
			Subtitle("Add more target folder to configuration").
			Valid(true).
			Var("action", "add")
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(t2fCmd)
}
