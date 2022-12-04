/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2here/alfred"
)

// t2hCmd represents the t2h command
var t2hCmd = &cobra.Command{
	Use:   "t2h",
	Short: "A brief description of your command",
	Run:   runT2hCmd,
}

func runT2hCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingTargets(wf)

	for name, path := range data {
		wi := wf.NewItem(name).
			Subtitle(fmt.Sprintf("↩ Move / Copy files from '%s' folder", path)).
			Valid(true).
			Arg(path)

		wi.Cmd().
			Subtitle("↩ Enter Action menu to add / remove target folder").
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
	rootCmd.AddCommand(t2hCmd)
}
