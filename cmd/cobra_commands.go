package cmd

import "github.com/spf13/cobra"

func WireUpCobraHarness(harness *CobraHarness) *cobra.Command {
	root := &cobra.Command{
		Use: "spec",
	}

	root.SetOutput(harness.stdout)

	root.AddCommand(
		commandConfig(harness),
		commandMetadata(harness),
	)

	root.PersistentPreRunE = harness.PersistentPreRunE

	return root
}

func commandConfig(harness *CobraHarness) *cobra.Command {
	root := &cobra.Command{
		Use: "config",
	}
	list := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
	}
	get := &cobra.Command{
		Use:     "get <key>",
		Args:    cobra.ExactArgs(1),
		Example: "$ spec config get project.name",
	}
	set := &cobra.Command{
		Use:     "set <key>=<value>",
		Example: "$ spec config set project.name=myProject",
	}

	root.AddCommand(
		list,
		get,
		set,
	)

	list.RunE = harness.ConfigList
	get.RunE = harness.ConfigGet
	set.Args = harness.ConfigSetArgs
	set.RunE = harness.ConfigSet

	return root
}

func commandMetadata(harness *CobraHarness) *cobra.Command {
	root := &cobra.Command{
		Use: "metadata",
	}
	add := &cobra.Command{
		Use:  "add",
		Args: cobra.MinimumNArgs(2),
	}
	root.AddCommand(
		add,
	)

	root.PersistentFlags().String("story", "", "")
	add.RunE = harness.MetadataAdd

	return root
}
