package logger

import (
	"bufio"
	"log"
	"os"

	"github.com/inconshreveable/mousetrap"
	"github.com/spf13/cobra"
)

var Logger = log.New(os.Stderr, "", 0)

func Exit(code int) {
	if mousetrap.StartedByExplorer() {
		Logger.Println("\nPress 'Enter' to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	os.Exit(code)
}

func init() {
	cobra.MousetrapHelpText = "" // Allow launching directly from Windows explorer
}
