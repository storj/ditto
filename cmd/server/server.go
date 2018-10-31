package server

import (
	minio "github.com/minio/minio/cmd"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	Args:    nil,
	Short:   "Upload files or file list_cmd to specified bucket",
	PreRunE: preRunE,
	Long:    `Upload files or file list_cmd to specified bucket`,
	Run:     run,
}

func preRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func run(cmd *cobra.Command, args []string) {
	minio.Main([]string{"mirroring", "gateway", "mirroring"})
}

func init() {

}
