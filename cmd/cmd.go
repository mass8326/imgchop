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

func Execute() {
	var flags RootFlags

	cmd := &cobra.Command{
		Use:   filepath.Base(os.Args[0]) + " [flags] [paths...]",
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
				go imgchop.Process(&wg, &warning, file, *flags.intelligent)
			}
			wg.Wait()

			if warning {
				logger.Exit(0)
			}
		},
		Version: version,
	}

	flags = initFlags(cmd)

	err := cmd.Execute()
	if err != nil {
		logger.Exit(1)
	}
}
