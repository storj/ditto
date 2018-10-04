package get

import "github.com/spf13/cobra"

func validateArgs(cmd *cobra.Command, args []string) error {
	argsLen := len(args)
	if argsLen < minArg || argsLen > maxArg {
		return NewArgsError(args)
	}

	return nil
}
