/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2go/alfred"
)

// t2glCmd represents the t2f command
var t2glCmd = &cobra.Command{
	Use:   "t2gl",
	Short: "Targets 2 go list",
	Run:   runT2glCmd,
}

func runT2glCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingTargets(wf)
	m, _ := cmd.Flags().GetString("mode")

	for name, path := range data {
		wi := wf.NewItem(name).
			Subtitle(fmt.Sprintf("⌘ ,↩ Move / Copy files from '%s' to frontmost Finder", path)).
			Valid(true).
			Arg(path).
			Var("mode", m)

		wi.Cmd().
			Subtitle("↩ Enter Action menu to Add / Remove target folder").
			Valid(true).
			Var("mode", m)
	}

	if len(data) == 0 {
		wf.NewItem("No any Targets available").
			Subtitle("↩ Add more target folder to configuration").
			Valid(true).
			Arg("").
			Var("mode", m)
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(t2glCmd)
	t2glCmd.PersistentFlags().StringP("mode", "m", "", "mode of action: to or go")
}
