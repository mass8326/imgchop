package logger

import (
	"bufio"
	"log"
	"os"

	"github.com/mass8326/imgchop/lib/util"
	"github.com/spf13/cobra"
)

var Logger = log.New(os.Stderr, "", 0)

func Exit(code int) {
	if util.StartedByExplorer {
		Logger.Println("\nPress 'Enter' to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	os.Exit(code)
}

func init() {
	cobra.MousetrapHelpText = "" // Allow launching directly from Windows explorer
}
