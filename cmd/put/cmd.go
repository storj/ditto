package put

import (
	"github.com/spf13/cobra"
	"os"
	"storj.io/ditto/cmd/utils"

	l "storj.io/ditto/pkg/logger"
)


// Cmd represents the put command
var Cmd = &cobra.Command{
	Use: "put [bucket name] [path to file/folder]",
	Args:    validateArgs,
	Aliases: []string{"p"},
	Short:   "Upload files or file list_cmd to specified bucket",
	Long:    `Upload files or file list_cmd to specified bucket`,
	RunE:    NewPutExec(utils.GetGateway, &l.StdOutLogger).runE,
}

var (
	frecursive, fforce bool
	fprefix, fdelimiter string // TODO: implement delimiter flag

	sigc chan os.Signal
)

func init() {
	//TODO: investigate on delayed response
	Cmd.Flags().BoolVarP(&frecursive, "recursive", "r", false, "recursive usage")
	Cmd.Flags().BoolVarP(&fforce, "force", "f", false, "force usage")
	Cmd.Flags().StringVarP(&fprefix, "prefix", "p", "", "prefix usage")
	Cmd.Flags().StringVarP(&fdelimiter, "delimiter", "d", "/", "delimiter usage")
}