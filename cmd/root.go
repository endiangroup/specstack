package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:           "specstack",
	SilenceErrors: true,
}

func init() {
	Config.AddCommand(ConfigList)

	Root.AddCommand(Config)
}

func WireUpHarness(harness *CobraHarness) {
	Root.SetOutput(harness.stdout)

	Root.PersistentPreRunE = harness.PersistentPreRunE
	ConfigList.RunE = harness.ConfigList
}
