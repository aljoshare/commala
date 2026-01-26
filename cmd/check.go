package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aljoshare/commala/internal/cli"
	"github.com/aljoshare/commala/internal/config"
	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/report"
	"github.com/aljoshare/commala/internal/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	MainCmd.AddCommand(queryCmd)
	queryCmd.Flags().String("config", "", "Path to the configuration file")
	viper.BindEnv("config", "COMMALA_CONFIG")
	viper.BindPFlag("config", queryCmd.Flags().Lookup("config"))
	queryCmd.Flags().String("report-junit-path", "commala-junit.xml", "Path of the JUnit report")
	viper.BindEnv("report.junit.path", "COMMALA_REPORT_JUNIT_PATH")
	viper.BindPFlag("report.junit.path", queryCmd.Flags().Lookup("report-junit-path"))
	queryCmd.Flags().Bool("author-email-enabled", true, "Flag to enable/disable author email validation")
	viper.BindEnv("validate.author.email.enabled", "COMMALA_VALIDATE_AUTHOR_EMAIL_ENABLED")
	viper.BindPFlag("validate.author.email.enabled", queryCmd.Flags().Lookup("author-email-enabled"))
	queryCmd.Flags().Bool("author-name-enabled", true, "Flag to enable/disable author name validation")
	viper.BindEnv("validate.author.name.enabled", "COMMALA_VALIDATE_AUTHOR_NAME_ENABLED")
	viper.BindPFlag("validate.author.name.enabled", queryCmd.Flags().Lookup("author-name-enabled"))
	queryCmd.Flags().Bool("branch-enabled", true, "Flag to enable/disable branch name validation")
	viper.BindEnv("validate.branch.enabled", "COMMALA_VALIDATE_BRANCH_ENABLED")
	viper.BindPFlag("validate.branch.enabled", queryCmd.Flags().Lookup("branch-enabled"))
	queryCmd.Flags().Bool("message-enabled", true, "Flag to enable/disable commit message validation")
	viper.BindEnv("validate.message.enabled", "COMMALA_VALIDATE_MESSAGE_ENABLED")
	viper.BindPFlag("validate.message.enabled", queryCmd.Flags().Lookup("message-enabled"))
	queryCmd.Flags().Bool("signoff-enabled", true, "Flag to enable/disable sign-off validation")
	viper.BindEnv("validate.signoff.enabled", "COMMALA_VALIDATE_SIGNOFF_ENABLED")
	viper.BindPFlag("validate.signoff.enabled", queryCmd.Flags().Lookup("signoff-enabled"))

	queryCmd.Flags().StringSlice("branch-whitelist", []string{}, "Contributor emails to whitelist for branch validation")
	viper.BindEnv("validate.branch.whitelist", "COMMALA_VALIDATE_BRANCH_WHITELIST")
	viper.BindPFlag("validate.branch.whitelist", queryCmd.Flags().Lookup("branch-whitelist"))

	queryCmd.Flags().StringSlice("message-whitelist", []string{}, "Contributor emails to whitelist for message validation")
	viper.BindEnv("validate.message.whitelist", "COMMALA_VALIDATE_MESSAGE_WHITELIST")
	viper.BindPFlag("validate.message.whitelist", queryCmd.Flags().Lookup("message-whitelist"))

	queryCmd.Flags().StringSlice("signoff-whitelist", []string{}, "Contributor emails to whitelist for signoff validation")
	viper.BindEnv("validate.signoff.whitelist", "COMMALA_VALIDATE_SIGNOFF_WHITELIST")
	viper.BindPFlag("validate.signoff.whitelist", queryCmd.Flags().Lookup("signoff-whitelist"))

	queryCmd.Flags().StringSlice("author-name-whitelist", []string{}, "Contributor emails to whitelist for author name validation")
	viper.BindEnv("validate.author.name.whitelist", "COMMALA_VALIDATE_AUTHOR_NAME_WHITELIST")
	viper.BindPFlag("validate.author.name.whitelist", queryCmd.Flags().Lookup("author-name-whitelist"))

	queryCmd.Flags().StringSlice("author-email-whitelist", []string{}, "Contributor emails to whitelist for author email validation")
	viper.BindEnv("validate.author.email.whitelist", "COMMALA_VALIDATE_AUTHOR_EMAIL_WHITELIST")
	viper.BindPFlag("validate.author.email.whitelist", queryCmd.Flags().Lookup("author-email-whitelist"))
}

var queryCmd = &cobra.Command{
	Use:   "check",
	Short: "Check commits",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Config{}
		c.ReadConfig()
		g := git.RealGit{}
		var cr *git.CommitRange
		var err error
		if isCommitRange(args[0]) {
			cr, err = g.ParseCommitRange(args[0])
		} else if isNegativeIndex(args[0]) {
			cr, err = g.ParseNegativeIndex(args[0])
		} else {
			cli.ErrorHandling(fmt.Errorf("argument must be a commit range or a negative index"))
			os.Exit(1)
		}
		if err != nil {
			cli.ErrorHandling(err)
			os.Exit(1)
		}
		result, err := validator.Validate(cr, g, c)
		if err != nil {
			cli.ErrorHandling(err)
			os.Exit(1)
		}
		cli.PrintResultTable(result)
		report.NewJUnitReport(result, c.ReportJunitPath)
		for _, r := range result {
			if !r.Valid {
				os.Exit(1)
			}
		}
	},
}

func isCommitRange(arg string) bool {
	return strings.Contains(arg, "..")
}

func isNegativeIndex(arg string) bool {
	return strings.HasPrefix(arg, "HEAD~")
}
