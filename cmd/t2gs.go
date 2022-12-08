/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-source2target/alfred"
)

// t2gsCmd represents the t2gs command
var t2gsCmd = &cobra.Command{
	Use:   "t2gs",
	Short: "Source 2 Target settings",
	Run:   runt2gsCmd,
}

func runt2gsCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingSources(wf)

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
	rootCmd.AddCommand(t2gsCmd)
}
