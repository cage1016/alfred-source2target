/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-source2target/alfred"
	"github.com/cage1016/alfred-source2target/lib"
)

// dotoCmd represents the doto command
var dogoCmd = &cobra.Command{
	Use:   "dogo",
	Short: "do go there",
	Run:   runDogoCmd,
}

func runDogoCmd(ccmd *cobra.Command, args []string) {
	source, _ := ccmd.Flags().GetString("source")

	documents := lib.FdExecute(lib.DoConfig{
		Source:    source,
		Arg:       args[0],
		Type:      alfred.GetType(wf),
		Exclude:   alfred.GetExclude(wf),
		ExecBatch: alfred.GetExecBatch(wf),
	})

	if len(documents) > 0 {
		if !strings.HasPrefix(documents[0], "error") {
			for i, doc := range documents {
				if strings.TrimSpace(doc) == "" {
					continue
				}

				fdoc := filepath.Join(source, doc)
				wf.NewItem(fmt.Sprintf("%d - %s", i, doc)).
					Subtitle("Press return to move recent files up to this file into target folder.").
					Quicklook(fdoc).
					Valid(true).
					Icon(&aw.Icon{
						Value: fdoc,
						Type:  "fileicon",
					}).
					ActionForType("file", fdoc).
					Arg(ArgJSONBuilder("move", documents[:i+1], source))

				// wi.Cmd().
				// 	Subtitle("Press return to move only this file into target folder.").
				// 	Valid(true).
				// 	Arg(ArgJSONBuilder("move", []string{doc}, target))

				// wi.Alt().
				// 	Subtitle("Press return to copy recent files up to this file into target folder.").
				// 	Valid(true).
				// 	Arg(ArgJSONBuilder("copy", documents[:i+1], target))

				// wi.NewModifier("alt", "cmd").
				// 	Subtitle("Press return to copy only this file into target folder.").
				// 	Valid(true).
				// 	Arg(ArgJSONBuilder("copy", []string{doc}, target))

				// wi.Shift().
				// 	Subtitle("Press return to send recent files up to this file to the pasteboard.").
				// 	Valid(true).
				// 	Arg(ArgJSONBuilder("pasteboard", documents[:i+1], target))

				// wi.NewModifier("cmd", "shift").
				// 	Subtitle("Press return to send only this file to the pasteboard.").
				// 	Valid(true).
				// 	Arg(ArgJSONBuilder("pasteboard", []string{doc}, target))

				// wi.Ctrl().
				// 	Subtitle("Browse in alfred").
				// 	Valid(true).
				// 	Arg(fdoc)
			}
		} else {
			wf.NewItem(documents[0]).Subtitle("").Valid(false)
		}
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(dogoCmd)
	dogoCmd.PersistentFlags().StringP("source", "", "", "source folder")
}
