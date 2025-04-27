package config

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type params struct {
	fx.In

	CobraCLI *cobra.Command
}

type result struct {
	fx.Out

	Config *Config
}

func newConfig(p params) (result, error) {
	service, err := newService(p.CobraCLI)

	return result{
		Config: service,
	}, err
}

var Module = fx.Module(
	"ConfigModule",
	fx.Provide(
		newConfig,
	),
)
