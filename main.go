package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/cmd"
)

func main() {
	var root = &cobra.Command{
		Use: "timelines",
	}
	root.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	cmd.Initialize(root)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
