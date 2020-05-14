package logging

import (
	"github.com/spf13/cobra"
)

// LevelFlagName is the canonical flag name to configure the allowed log level
// within Prometheus projects.
const LevelFlagName = "log.level"

// LevelFlagHelp is the help description for the log.level flag.
const LevelFlagHelp = "Only log messages with the given severity or above. One of: [debug, info, warn, error]"

// FormatFlagName is the canonical flag name to configure the log format
// within Prometheus projects.
const FormatFlagName = "log.format"

// FormatFlagHelp is the help description for the log.format flag.
const FormatFlagHelp = "Output format of log messages. One of: [logfmt, json]"

// AddFlags adds the flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddFlags(command *cobra.Command, config *Config) {
	config.Level = &AllowedLevel{}
	config.Level.Set("info")
	command.PersistentFlags().Var(config.Level, LevelFlagName, LevelFlagHelp)

	config.Format = &AllowedFormat{}
	config.Format.Set("logfmt")
	command.PersistentFlags().Var(config.Format, FormatFlagName, FormatFlagHelp)
}
