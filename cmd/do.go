/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2go/alfred"
	"github.com/cage1016/alfred-targets2go/lib"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "A brief description of your command",
	Run:   runDoCmd,
}

type Arg struct {
	Op    string   `json:"op"`
	Files []string `json:"files"`
	Base  string   `json:"base"`
}

func ArgJSONBuilder(op string, files []string, base string) string {
	j, _ := json.Marshal(Arg{op, files, base})
	return string(j)
}

func runDoCmd(ccmd *cobra.Command, args []string) {
	target, _ := ccmd.Flags().GetString("target")

	documents := lib.FdExecute(lib.DoConfig{
		Target:    target,
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

				fdoc := filepath.Join(target, doc)
				wi := wf.NewItem(fmt.Sprintf("%d - %s", i, doc)).
					Subtitle("Press return to move recent files up to this file into the current Finder location.").
					Quicklook(fdoc).
					Valid(true).
					Icon(&aw.Icon{
						Value: fdoc,
						Type:  "fileicon",
					}).
					ActionForType("file", fdoc).
					Arg(ArgJSONBuilder("move", documents[:i+1], target))

				wi.Cmd().
					Subtitle("Press return to move only this file into the current Finder location.").
					Valid(true).
					Arg(ArgJSONBuilder("move", []string{doc}, target))

				wi.Alt().
					Subtitle("Press return to copy recent files up to this file into the current Finder location.").
					Valid(true).
					Arg(ArgJSONBuilder("copy", documents[:i+1], target))

				wi.NewModifier("alt", "cmd").
					Subtitle("Press return to copy only this file into the current Finder location.").
					Valid(true).
					Arg(ArgJSONBuilder("copy", []string{doc}, target))

				wi.Shift().
					Subtitle("Press return to send recent files up to this file to the pasteboard.").
					Valid(true).
					Arg(ArgJSONBuilder("pasteboard", documents[:i+1], target))

				wi.NewModifier("cmd", "shift").
					Subtitle("Press return to send only this file to the pasteboard.").
					Valid(true).
					Arg(ArgJSONBuilder("pasteboard", []string{doc}, target))

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
	rootCmd.AddCommand(doCmd)
	doCmd.PersistentFlags().StringP("target", "", "", "target folder")
}
