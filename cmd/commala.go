package cmd

import (
	"os"

	"github.com/aljoshare/commala/internal/logging"
	"github.com/spf13/cobra"
)

func init() {

}

var MainCmd = &cobra.Command{
	Use:   "commala",
	Short: "commala is a commit linter with a lot of rice",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	log, _ := logging.GetLogger()
	if err := MainCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
