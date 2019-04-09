package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/cmd"
)

func main() {
	var envFiles []string
	root := &cobra.Command{
		Use: "timelines",
		PersistentPreRun: func(*cobra.Command, []string) {
			godotenv.Load()
			if len(envFiles) > 0 {
				godotenv.Load(envFiles...)
			}
		},
	}
	root.Flags().StringArrayVar(&envFiles, "env", nil, "env files to load")
	// hide the help command because it's not pretty
	root.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})

	// initialize and execute commands
	cmd.Initialize(root)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
