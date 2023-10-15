/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-source2target/alfred"
	"github.com/cage1016/alfred-source2target/lib"
)

// docomeCmd represents the docome command
var docomeCmd = &cobra.Command{
	Use:   "docome",
	Short: "do come here",
	Run:   runDocomeCmd,
}

func runDocomeCmd(ccmd *cobra.Command, args []string) {
	source, _ := ccmd.Flags().GetString("source")

	var err error
	ranges := &[]lib.Range{}
	query := args[0]
	if lib.IsPageRangeValid(args[0]) {
		page := lib.PageRangeRegex.FindAllString(args[0], -1)
		query = strings.TrimSpace(strings.Replace(args[0], page[0], "", -1))
		ranges, err = lib.ParseRangeNumber(strings.Replace(page[0], "#", "", -1), alfred.GetMaxQueryResults(wf))
		if err != nil {
			ErrorHandle(fmt.Errorf("failed to parse page range: %v", err))
			return
		}
	}

	documents := lib.FdExecute(lib.DoConfig{
		Source:          source,
		Arg:             query,
		Type:            alfred.GetType(wf),
		Exclude:         alfred.GetExclude(wf),
		MaxQueryResults: alfred.GetMaxQueryResults(wf),
	})

	// filter targets for page range case
	nTargets := lo.Filter(documents, func(_ string, index int) bool {
		return lib.Ranges(*ranges).IsInRange(index + 1)
	})

	if len(documents) > 0 {
		if !lib.IsFdError(documents[0]) {
			for i, doc := range documents {
				if strings.TrimSpace(doc) == "" {
					continue
				}
				i += 1

				if len(*ranges) == 0 {
					fdoc := filepath.Join(source, doc)
					wi := wf.NewItem(fmt.Sprintf("%d - %s", i, doc)).
						Subtitle("⇧ ⌥ ⌘, Press return to move recent files up to this file into the current Finder location.").
						Quicklook(fdoc).
						Valid(true).
						Arg(fdoc).
						Icon(&aw.Icon{
							Value: fdoc,
							Type:  "fileicon",
						}).
						ActionForType("file", fdoc).
						Var("documents", ArgJSONBuilder("move", documents[:i], source)).
						Var("action", "")

					wi.Cmd().
						Subtitle("Press return to move only this file into the current Finder location.").
						Valid(true).
						Arg(fdoc).
						Var("documents", ArgJSONBuilder("move", []string{doc}, source)).
						Var("action", "")

					wi.Alt().
						Subtitle("Press return to copy recent files up to this file into the current Finder location.").
						Valid(true).
						Arg(fdoc).
						Var("documents", ArgJSONBuilder("copy", documents[:i], source)).
						Var("action", "")

					wi.NewModifier("alt", "cmd").
						Subtitle("Press return to copy only this file into the current Finder location.").
						Valid(true).
						Arg(fdoc).
						Var("documents", ArgJSONBuilder("copy", []string{doc}, source)).
						Var("action", "")

					wi.Shift().
						Subtitle("Press return to send recent files up to this file to the pasteboard.").
						Valid(true).
						Arg(fdoc).
						Var("documents", ArgJSONBuilder("pasteboard", documents[:i], source)).
						Var("action", "")

					wi.NewModifier("cmd", "shift").
						Subtitle("Press return to send only this file to the pasteboard.").
						Valid(true).
						Arg(fdoc).
						Var("documents", ArgJSONBuilder("pasteboard", []string{doc}, source)).
						Var("action", "")

					wi.Ctrl().
						Subtitle("Browse in alfred").
						Valid(true).
						Arg(fdoc).
						Var("action", "")
				} else {
					fdoc := filepath.Join(source, doc)
					var wi *aw.Item
					if lib.Ranges(*ranges).IsInRange(i) {
						wi = wf.NewItem(fmt.Sprintf("✅ %d - %s", i, doc)).
							Subtitle("Press return to choose action into the current Finder location.")
					} else {
						wi = wf.NewItem(fmt.Sprintf("%d - %s", i, doc))
					}

					wi.Quicklook(fdoc).
						Valid(true).
						Arg(fdoc).
						Icon(&aw.Icon{
							Value: fdoc,
							Type:  "fileicon",
						}).
						ActionForType("file", fdoc).
						Var("documents", ArgJSONBuilder("", nTargets, source)).
						Var("action", "choose op")
				}
			}
		} else {
			for _, doc := range documents {
				wf.NewItem(doc).
					Subtitle("Press return to visit fd document").
					Valid(true).
					Var("action", "open").
					Var("help", "https://github.com/sharkdp/fd")
			}
		}
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(docomeCmd)
	docomeCmd.PersistentFlags().StringP("source", "", "", "source folder")
}
