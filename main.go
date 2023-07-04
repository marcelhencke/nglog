package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"nglog/config"
	debug "nglog/debugPrint"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "nglog [flags] [LOG_FILE]",
	Short: "Formats (php-fpm + nginx) logs",
	Long:  `nglog read text lines from LOG_FILE or from PIPE. This lines gets buffered until the LogLineCompleteRegex matches. After this nglog trys to parse a typical php error line in a nginx error log file (LogLineFastCGIRegex). If this parsing is successful, nglog uses this data to print a new log line defined by the template. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			config.Flags.FilePath = args[0]
		}
		debug.Println = debug.PrintNoop
		if config.Flags.DebugMode {
			debug.Println = debug.PrintFmt
		}
		return LoadData()
	},
}

func main() {
	rootCmd.Flags().StringVarP(
		&config.Flags.Template,
		"template",
		"t",
		"%raw%",
		"This is a template format for the new line. The key must be enclosed in %\n"+
			"Possible keys:\n"+
			"- raw: no further line manipulation\n"+
			"- prefix: first match group from LogLineFastCGIRegex\n"+
			"- suffix: third match group from LogLineFastCGIRegex\n"+
			"- ts: timestamp\n"+
			"- php: PHP message\n"+
			"- ng_xxx: nginx var, e.g. ng_server, ng_upstream ...")

	rootCmd.Flags().StringVar(
		&config.Flags.LogLineCompleteRegex,
		"overwriteLogLineCompleteRegex",
		"\" while reading[A-z ]* upstream",
		"This RegEx tests whether the lines read are a complete log line.\n")
	rootCmd.Flags().StringVar(
		&config.Flags.LogLineFastCGIRegex,
		"overwriteLogLineFastCGIRegex",
		"^(\\d{4}[\\/-]\\d{2}[\\/-]\\d{2} \\d{2}:\\d{2}:\\d{2} \\[\\w+\\] .+ FastCGI sent in stderr: \")([\\s\\S]*)(\" while reading[A-z ]* upstream(?:, \\w+: \"?.+\"?)*)$",
		"This RegEx tests whether the log line is a FastCGI log line.\nThis expression are divided in 3 matching groups: prefix, body, suffix.\nPrefix: Timestamp and FastCGI identifier.\nBody: Contains the PHP messages.\nSuffix: Some nginx properties, like client, server, upstream, host, etc.\n")
	rootCmd.Flags().StringVar(
		&config.Flags.LogLinePhpMsgSplit,
		"overwriteLogLinePhpMsgSplit",
		"PHP message: ",
		"This string is the split value for PHP messages.\n")
	rootCmd.Flags().StringVar(
		&config.Flags.NginxVarRegex,
		"overwriteNginxVarRegex",
		", (\\w+): (\"?[^,]+\"?)",
		"This RegEx finds the nginx vars in the line suffix.\n")

	rootCmd.PersistentFlags().BoolVarP(
		&config.Flags.DebugMode,
		"debugMode",
		"d",
		false, "Print out debug logs. Helpful when defining custom regex.")
	rootCmd.Flags().IntVar(
		&config.Flags.DebugMaxLines,
		"debugMaxLines",
		0,
		"Only read x lines of input data. Only effective when debugMode is enabled and x > 0. ")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error message: " + err.Error() + ".")
		os.Exit(1)
	}
}
