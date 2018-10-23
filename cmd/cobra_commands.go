package cmd

import "github.com/spf13/cobra"

var (
	cmdRoot = &cobra.Command{
		Use: "spec",
	}

	cmdConfig = &cobra.Command{
		Use: "config",
	}

	cmdConfigList = &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
	}
	cmdConfigGet = &cobra.Command{
		Use:     "get <key>",
		Args:    cobra.MinimumNArgs(1),
		Example: "$ spec config get project.name",
	}
	cmdConfigSet = &cobra.Command{
		Use:     "set <key>=<value>",
		Example: "$ spec config set project.name=myProject",
	}
)

func WireUpCobraHarness(harness *CobraHarness) *cobra.Command {
	cmdConfig.AddCommand(cmdConfigList)
	cmdConfig.AddCommand(cmdConfigGet)
	cmdConfig.AddCommand(cmdConfigSet)
	cmdRoot.AddCommand(cmdConfig)

	cmdRoot.SetOutput(harness.stdout)

	cmdRoot.PersistentPreRunE = harness.PersistentPreRunE
	cmdConfigList.RunE = harness.ConfigList
	cmdConfigGet.RunE = harness.ConfigGet
	cmdConfigSet.Args = harness.ConfigSetArgs
	cmdConfigSet.RunE = harness.ConfigSet

	return cmdRoot
}
