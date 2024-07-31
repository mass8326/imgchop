package logger

import (
	"bufio"
	"log"
	"os"

	"github.com/mass8326/imgchop/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var RawLogger = log.New(os.Stderr, "", 0)

type Message struct {
	Level  string
	Source string
	Msg    string
}

type Logger struct {
	Messages chan Message
	Source   string
}

func (lg Logger) Warn(msg string) {
	lg.Messages <- Message{
		Level:  "warn",
		Source: lg.Source,
		Msg:    msg,
	}
}

func (lg Logger) Info(msg string) {
	lg.Messages <- Message{
		Level:  "info",
		Source: lg.Source,
		Msg:    msg,
	}
}

func Exit(code int) {
	if util.StartedByExplorer {
		previous, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			RawLogger.Println("\nPress 'Enter' to exit...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		} else {
			defer term.Restore(int(os.Stdin.Fd()), previous)
			RawLogger.Println("\nPress any key to exit...")
			b := make([]byte, 1)
			os.Stdin.Read(b)
		}
	}
	os.Exit(code)
}

func init() {
	cobra.MousetrapHelpText = "" // Allow launching directly from Windows explorer
}
