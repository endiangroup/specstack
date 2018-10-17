package cmd

import "github.com/spf13/cobra"

var (
	cmdRoot = &cobra.Command{
		Use:           "specstack",
		SilenceErrors: true,
	}

	cmdConfig = &cobra.Command{
		Use: "config",
	}

	cmdConfigList = &cobra.Command{
		Use: "list",
	}
)

func WireUpCobraHarness(harness *CobraHarness) *cobra.Command {
	cmdConfig.AddCommand(cmdConfigList)
	cmdRoot.AddCommand(cmdConfig)

	cmdRoot.SetOutput(harness.stdout)

	cmdRoot.PersistentPreRunE = harness.PersistentPreRunE
	cmdConfigList.RunE = harness.ConfigList

	return cmdRoot
}
