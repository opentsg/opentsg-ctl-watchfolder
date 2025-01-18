package cli

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/opentsg/opentsg-ctl-watchfolder/semver"
	"github.com/spf13/cobra"
)

// mostly used for the version string
const App = "opentsg-ctl-watchfolder"

// these variables are set by the cobra CLI handler (see cli-main.init())
var ShowVersion bool
var ShowVersionShort bool
var ShowVersionNote bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version string",
	Long:  `use -v for short semver and --version for long seMver.`,
	Run: func(cmd *cobra.Command, args []string) {

		if ShowVersionShort {
			fmt.Println(semver.Info.Short)

		} else if ShowVersionNote {
			fmt.Println(semver.Info.Note)

		} else {
			fmt.Printf("%s %s\n", App, semver.Info.Long)
		}
	},
}

func init() {
	// trace init order for sanity
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	// add this command to the main command
	mainCmd.AddCommand(versionCmd)
}
