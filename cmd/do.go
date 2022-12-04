/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

func worker(done chan struct{}, r io.ReadCloser, fn func(string)) {
	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			fn(scanner.Text())
		}
		done <- struct{}{}
	}()
}

func runDoCmd(ccmd *cobra.Command, args []string) {
	target, _ := ccmd.Flags().GetString("target")

	flags := []string{}
	for _, arg := range strings.Split(args[0], " ") {
		if arg != "" {
			flags = append(flags, arg)
		}
	}
	flags = append(flags, "-t", "f", "-X", "ls", "-lt")

	cmd := exec.Command("fd", flags...)
	cmd.Dir = target
	logrus.Debugf("fd: %s", cmd)

	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	documents := []string{}
	done := make(chan struct{})
	worker(done, r, func(line string) {
		logrus.Debugf("line: %s", line)
		if strings.HasPrefix(line, "-") {
			x := strings.Split(line, "./")
			documents = append(documents, x[1])
		} else {
			documents = append(documents, line)
		}
	})

	cmd.Start()
	<-done
	cmd.Wait()

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
