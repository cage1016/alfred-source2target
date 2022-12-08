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
				wi := wf.NewItem(fmt.Sprintf("%d - %s", i, doc)).
					Subtitle("Press return to move recent files up to this file into target folder.").
					Valid(true).
					Arg("").
					Icon(&aw.Icon{Value: fdoc, Type: "fileicon"}).
					Quicklook(fdoc).
					ActionForType("file", fdoc).
					Var("query", ArgJSONBuilder("move", documents[:i+1], source))

				wi.Cmd().
					Subtitle("Press return to move only this file into target folder.").
					Valid(true).
					Arg("").
					Var("query", ArgJSONBuilder("move", []string{doc}, source))

				wi.Alt().
					Subtitle("Press return to copy recent files up to this file into target folder.").
					Valid(true).
					Arg("").
					Var("query", ArgJSONBuilder("copy", documents[:i+1], source))

				wi.NewModifier("alt", "cmd").
					Subtitle("Press return to copy only this file into target folder.").
					Valid(true).
					Arg("").
					Var("query", ArgJSONBuilder("copy", []string{doc}, source))

				wi.Shift().
					Subtitle("Press return to send recent files up to this file to the pasteboard.").
					Valid(true).
					Arg("").
					Var("query", ArgJSONBuilder("pasteboard", documents[:i+1], source))

				wi.NewModifier("cmd", "shift").
					Subtitle("Press return to send only this file to the pasteboard.").
					Valid(true).
					Arg("").
					Var("query", ArgJSONBuilder("pasteboard", []string{doc}, source))

				wi.Ctrl().
					Subtitle("Browse in alfred").
					Valid(true).
					Arg(fdoc)
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
