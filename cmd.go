package notes

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

// Cmd is an interface for subcommands of notes command
type Cmd interface {
	Do() error
	defineCLI(*kingpin.Application)
	matchesCmdline(string) bool
}

// Version is version string of notes command
var Version = "1.0.0"
