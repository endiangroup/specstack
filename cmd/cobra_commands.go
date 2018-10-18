package cmd

import "github.com/spf13/cobra"

var (
	cmdRoot = &cobra.Command{
		Use:           "spec",
		SilenceErrors: true,
	}

	cmdConfig = &cobra.Command{
		Use: "config",
	}

	cmdConfigList = &cobra.Command{
		Use: "list",
	}
	cmdConfigGet = &cobra.Command{
		Use:  "get",
		Args: cobra.MinimumNArgs(1),
	}
)

func WireUpCobraHarness(harness *CobraHarness) *cobra.Command {
	cmdConfig.AddCommand(cmdConfigList)
	cmdConfig.AddCommand(cmdConfigGet)
	cmdRoot.AddCommand(cmdConfig)

	cmdRoot.SetOutput(harness.stdout)

	cmdRoot.PersistentPreRunE = harness.PersistentPreRunE
	cmdConfigList.RunE = harness.ConfigList
	cmdConfigGet.RunE = harness.ConfigGet

	return cmdRoot
}
