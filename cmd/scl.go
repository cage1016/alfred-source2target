/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-source2target/alfred"
)

// sclCmd represents the scl command
var sclCmd = &cobra.Command{
	Use:   "scl",
	Short: "Source come list",
	Run:   runSclCmd,
}

func runSclCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingSources(wf)
	m, _ := cmd.Flags().GetString("mode")

	for name, path := range data {
		wi := wf.NewItem(name).
			Subtitle(fmt.Sprintf("⌘ ,↩ Move / Copy files from '%s' to frontmost Finder", path)).
			Valid(true).
			Arg(path).
			Var("mode", m)

		wi.Cmd().
			Subtitle("↩ Enter Action menu to Add / Remove source folder").
			Valid(true).
			Arg("").
			Var("mode", m)
	}

	if len(data) == 0 {
		wf.NewItem("No any source available").
			Subtitle("↩ Add more source folder to configuration").
			Valid(true).
			Arg("").
			Var("mode", m)
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(sclCmd)
	sclCmd.PersistentFlags().StringP("mode", "m", "", "mode of action: to or go")
}
