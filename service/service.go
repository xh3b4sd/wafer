// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/spf13/viper"

	"github.com/xh3b4sd/wafer/flag"
	"github.com/xh3b4sd/wafer/service/analyzer"
	v1analyzer "github.com/xh3b4sd/wafer/service/analyzer/v1"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/informer/csv"
	"github.com/xh3b4sd/wafer/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var informerService informer.Informer
	{
		informerConfig := csv.DefaultConfig()
		//		informerConfig.File.Header.Buy = 1
		//		informerConfig.File.Header.Ignore = false
		//		informerConfig.File.Header.Sell = 2
		//		informerConfig.File.Header.Time = 0
		//		informerConfig.File.Path = "/Users/xh3b4sd/go/src/github.com/xh3b4sd/wafer/charts/001/chart.csv"
		informerConfig.Dir.Path = "/Users/xh3b4sd/go/src/github.com/xh3b4sd/wafer/charts/"
		informerService, err = csv.New(informerConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var analyzerService analyzer.Analyzer
	{
		analyzerConfig := v1analyzer.DefaultConfig()
		analyzerConfig.Informer = informerService
		analyzerConfig.Logger = config.Logger
		analyzerService, err = v1analyzer.New(analyzerConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Logger = config.Logger

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		Analyzer: analyzerService,
		Version:  versionService,
	}

	return newService, nil
}

type Service struct {
	Analyzer analyzer.Analyzer
	Version  *version.Service
}
