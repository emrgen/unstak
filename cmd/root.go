package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "unp",
	Short: "unpost CLI",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(spaceCommand)
	rootCmd.AddCommand(tierCmd)
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(tagCmd)
	rootCmd.AddCommand(collectionCmd)
}
