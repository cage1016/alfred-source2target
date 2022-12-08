/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-source2target/alfred"
)

// sglCmd represents the sgl command
var sglCmd = &cobra.Command{
	Use:   "sgl",
	Short: "Source go list",
	Run:   runSglCmd,
}

func runSglCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingSources(wf)
	m, _ := cmd.Flags().GetString("mode")

	for name, path := range data {
		wi := wf.NewItem(name).
			Subtitle(fmt.Sprintf("⌘ ,↩ Move / Copy files from '%s' to Target folder", path)).
			Valid(true).
			Arg(path).
			Var("mode", m)

		wi.Cmd().
			Subtitle("↩ Enter Action menu to Add / Remove source folder").
			Valid(true).
			Var("mode", m)
	}

	if len(data) == 0 {
		wf.NewItem("No any Source available").
			Subtitle("↩ Add more source folder to configuration").
			Valid(true).
			Arg("").
			Var("mode", m)
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(sglCmd)
	sglCmd.PersistentFlags().StringP("mode", "m", "", "mode of action: to or go")
}
