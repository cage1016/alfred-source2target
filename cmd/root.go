/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2here/alfred"
)

const updateJobName = "checkForUpdate"

var (
	repo = "cage1016/alfred-targets2here"
	wf   *aw.Workflow
	av   = aw.NewArgVars()
)

func ErrorHandle(err error) {
	av.Var("error", err.Error())
	if err := av.Send(); err != nil {
		wf.Fatalf("failed to send args to Alfred: %v", err)
	}
}

func CheckForUpdate() {
	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		logrus.Info("Running update check in background...")
		cmd := exec.Command(os.Args[0], "update")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			logrus.Errorf("Error starting update check: %s", err)
		}
	}

	if wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		wf.NewItem("An update is available!").
			Subtitle("⇥ or ↩ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(&aw.Icon{Value: "update-available.png"})
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Targets 2 Here",
	Short: "Copy or Move targets to Here",
	Run: func(cmd *cobra.Command, args []string) {
		CheckForUpdate()
		wf.SendFeedback()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	wf.Run(func() {
		if _, err := os.Stat(filepath.Join(wf.DataDir(), "targets.json")); errors.Is(err, os.ErrNotExist) {
			alfred.StoreOngoingTargets(wf, alfred.Targets{})
		}

		if err := rootCmd.Execute(); err != nil {
			logrus.Fatal(err)
		}
	})
}

func init() {
	wf = aw.New(update.GitHub(repo), aw.HelpURL(repo+"/issues"))
	wf.Args() // magic for "workflow:update"

	if alfred.GetDebug(wf) {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
