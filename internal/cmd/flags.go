package cmd

import "github.com/spf13/cobra"

type RootFlags struct {
	intelligent *bool
}

func initFlags(cmd *cobra.Command) RootFlags {
	cmd.PersistentFlags().Bool("help", false, "print help message")
	cmd.PersistentFlags().Bool("version", false, "print version")
	cmd.SetVersionTemplate("{{.Version}}\n")

	return RootFlags{
		intelligent: cmd.PersistentFlags().BoolP("intelligent", "i", false, "enable intelligent filter"),
	}
}
