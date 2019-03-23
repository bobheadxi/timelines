package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/cmd"
)

func main() {
	var (
		envFiles []string
	)
	root := &cobra.Command{
		Use: "timelines",
		PersistentPreRun: func(*cobra.Command, []string) {
			godotenv.Load()
			if len(envFiles) > 0 {
				godotenv.Load(envFiles...)
			}
		},
	}
	root.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	root.Flags().StringArrayVar(&envFiles, "env", nil, "")
	cmd.Initialize(root)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
