/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-targets2go/alfred"
)

// maCmd represents the ma command
var maCmd = &cobra.Command{
	Use:   "ma",
	Short: "Manage Targets",
	Run:   runMaCmd,
}

func runMaCmd(cmd *cobra.Command, args []string) {
	data, _ := alfred.LoadOngoingTargets(wf)
	action, _ := cmd.Flags().GetString("action")
	folder := filepath.Base(args[0])

	switch action {
	case "add":
		data[folder] = args[0]
		alfred.StoreOngoingTargets(wf, data)
	case "remove":
		delete(data, folder)
		alfred.StoreOngoingTargets(wf, data)
	default:
		logrus.Info("Do nothing")
	}
}

func init() {
	rootCmd.AddCommand(maCmd)
	maCmd.PersistentFlags().StringP("action", "a", "", "type of Action")
}
