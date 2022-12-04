/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2here/alfred"
)

// t2hsCmd represents the t2hs command
var t2hsCmd = &cobra.Command{
	Use:   "t2hs",
	Short: "A brief description of your command",
	Run:   runt2hsCmd,
}

func runt2hsCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingTargets(wf)

	wf.NewItem("BACK").
		Subtitle("Back to list").
		Valid(true).
		Var("action", "back")

	wf.NewItem("ADD").
		Subtitle("Add more target folder to configuration").
		Valid(true).
		Var("action", "add")

	for name, path := range data {
		wf.NewItem(fmt.Sprintf("Remove %s", name)).
			Subtitle(fmt.Sprintf("remove '%s' from configuration", path)).
			Valid(true).
			Arg(path).
			Var("action", "remove")
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(t2hsCmd)
}
