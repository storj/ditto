package put

import (
	"github.com/spf13/cobra"
	"os"
	"storj.io/ditto/cmd/utils"

	l "storj.io/ditto/pkg/logger"
)


// Cmd represents the put command
var Cmd = &cobra.Command{
	Use: "put [bucket name] [path to file/folder] [OPTIONS]",
	Args:    validateArgs,
	Aliases: []string{"p"},
	Short:   "Upload file or folder to specified bucket",
	Long:    `Upload file or folder to specified bucket`,
	RunE:    NewPutExec(utils.GetGateway, &l.StdOutLogger).runE,
}

var (
	frecursive, fforce bool
	fprefix, fdelimiter string

	sigc chan os.Signal
)

func init() {
	//TODO: investigate on delayed response
	Cmd.Flags().BoolVarP(&frecursive, "recursive", "r", false, "recursively upload contents of the specified folder")
	Cmd.Flags().BoolVarP(&fforce, "force", "f", false, "truncate object if one exists")
	Cmd.Flags().StringVarP(&fprefix, "prefix", "p", "", "root prefix")
	Cmd.Flags().StringVarP(&fdelimiter, "delimiter", "d", "/", "separates objnames from prefixes")
}