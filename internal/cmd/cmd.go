package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/mass8326/imgchop/internal/imgchop"
	"github.com/mass8326/imgchop/internal/logger"
	"github.com/spf13/cobra"
)

// This should be defined at build time using ldflags
var version = "0.0.0-invalid"

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

			c := make(chan logger.Message, 1)
			var wg sync.WaitGroup
			wg.Add(len(args))
			for _, file := range args {
				go imgchop.Process(&wg, c, file, *flags.intelligent)
			}

			messaged := false
			go func() {
				for msg := range c {
					messaged = true
					var level string
					switch msg.Level {
					case "warn":
						level = color.New(color.FgYellow).Sprintf("[%s]", msg.Level)
					case "info":
						level = color.New(color.FgCyan).Sprintf("[%s]", msg.Level)
					}
					logger.RawLogger.Printf("%s %s", level, color.HiBlackString(msg.Source))
					logger.RawLogger.Println(msg.Msg)
				}
			}()

			wg.Wait()
			if messaged {
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
