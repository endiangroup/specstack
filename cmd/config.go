package cmd

import (
	"github.com/spf13/cobra"
)

var Config = &cobra.Command{
	Use: "config",
}

var ConfigList = &cobra.Command{
	Use: "list",
}
