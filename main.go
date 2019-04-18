package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/cmd"
	"github.com/bobheadxi/timelines/config"
)

func main() {
	var envFiles []string
	root := &cobra.Command{
		Use:     "timelines",
		Version: config.NewBuildMeta().Commit,
		PersistentPreRun: func(*cobra.Command, []string) {
			godotenv.Load()
			if len(envFiles) > 0 {
				godotenv.Load(envFiles...)
			}
		},
	}
	root.PersistentFlags().StringArrayVar(&envFiles, "env", nil, "env files to load")
	// hide the help command because it's not pretty
	root.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})

	// initialize and execute commands
	cmd.Initialize(root)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
