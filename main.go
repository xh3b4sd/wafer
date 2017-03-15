package main

import (
	"os"

	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/command"
)

var (
	description string = "Command line tool for automatically ride financial waves when stock market charts surge."
	gitCommit   string = "n/a"
	name        string = "wafer"
	source      string = "https://github.com/xh3b4sd/wafer"
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		newLogger, err = micrologger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	var newCommand *command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
		if err != nil {
			panic(err)
		}
	}

	newCommand.CobraCommand().Execute()
}
