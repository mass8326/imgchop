package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/mass8326/imgchop/lib/imgchop"
	"github.com/mass8326/imgchop/lib/logger"
	"github.com/spf13/cobra"
)

var version = "[N/A]"

var cmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]) + " [flags] [files...]",
	Short: "Chop your images into squares",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			logger.Exit(1)
		}

		warning := false
		var wg sync.WaitGroup
		for _, file := range args {
			wg.Add(1)
			go imgchop.Crop(&wg, &warning, file)
		}
		wg.Wait()

		if warning {
			logger.Exit(0)
		}
	},
	Version: version,
}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		logger.Exit(1)
	}
}

func init() {
	// These flags are automatically handled by Cobra
	cmd.PersistentFlags().Bool("help", false, "print help message")
	cmd.PersistentFlags().Bool("version", false, "print version")
	cmd.SetVersionTemplate("{{.Version}}\n")
}
