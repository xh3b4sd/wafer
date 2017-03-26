package main

import (
	"os"

	"github.com/giantswarm/microkit/command"
	"github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/spf13/viper"

	"github.com/xh3b4sd/wafer/flag"
	"github.com/xh3b4sd/wafer/server"
	"github.com/xh3b4sd/wafer/service"
)

var (
	description string     = "Microservice for automatically ride financial waves when stock market charts surge."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "wafer"
	source      string     = "https://github.com/xh3b4sd/wafer"
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger logger.Logger
	{
		loggerConfig := logger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		newLogger, err = logger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.Flag = f
			serviceConfig.Logger = newLogger
			serviceConfig.Viper = v

			serviceConfig.Description = description
			serviceConfig.GitCommit = gitCommit
			serviceConfig.Name = name
			serviceConfig.Source = source

			newService, err = service.New(serviceConfig)
			if err != nil {
				panic(err)
			}
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			serverConfig := server.DefaultConfig()

			serverConfig.MicroServerConfig.Logger = newLogger
			serverConfig.MicroServerConfig.ServiceName = name
			serverConfig.MicroServerConfig.Viper = v
			serverConfig.Service = newService

			newServer, err = server.New(serverConfig)
			if err != nil {
				panic(err)
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger
		commandConfig.ServerFactory = newServerFactory

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
		if err != nil {
			panic(err)
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Informer.CSV.File, "", "The absolute file path of a CSV file containing chart data.")
	daemonCommand.PersistentFlags().Bool(f.Service.Informer.CSV.IgnoreHeader, false, "Whether to ignore the first row of the CSV file.")
	daemonCommand.PersistentFlags().Int(f.Service.Informer.CSV.Index.Buy, 0, "The index of the column within a CSV file representing buy prices.")
	daemonCommand.PersistentFlags().Int(f.Service.Informer.CSV.Index.Sell, 0, "The index of the column within a CSV file representing sell prices.")
	daemonCommand.PersistentFlags().Int(f.Service.Informer.CSV.Index.Time, 0, "The index of the column within a CSV file representing price times.")
	daemonCommand.PersistentFlags().String(f.Service.Informer.Kind, "csv", "The kind of the informer imlementation to use.")

	newCommand.CobraCommand().Execute()
}
