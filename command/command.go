// Package command implements the root command for the command line tool.
package command

import (
	"github.com/spf13/cobra"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/command/render"
	"github.com/xh3b4sd/wafer/command/version"
)

// Config represents the configuration used to create a new root command.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new root command
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new root command.
func New(config Config) (*Command, error) {
	var err error

	var renderCommand *render.Command
	{
		renderConfig := render.DefaultConfig()

		renderConfig.Logger = config.Logger

		renderCommand, err = render.New(renderConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var versionCommand *version.Command
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionCommand, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newCommand := &Command{
		// Internals.
		cobraCommand:   nil,
		renderCommand:  renderCommand,
		versionCommand: versionCommand,
	}

	newCommand.cobraCommand = &cobra.Command{
		Use:   config.Name,
		Short: config.Description,
		Long:  config.Description,
		Run:   newCommand.Execute,
	}

	newCommand.cobraCommand.AddCommand(newCommand.renderCommand.CobraCommand())
	newCommand.cobraCommand.AddCommand(newCommand.versionCommand.CobraCommand())

	return newCommand, nil
}

type Command struct {
	// Internals.
	cobraCommand   *cobra.Command
	renderCommand  *render.Command
	versionCommand *version.Command
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) RenderCommand() *render.Command {
	return c.renderCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}

func (c *Command) VersionCommand() *version.Command {
	return c.versionCommand
}
