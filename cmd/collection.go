package cmd

import "github.com/spf13/cobra"

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "collection commands",
}

func init() {
	collectionCmd.AddCommand(collectionCreate())
	collectionCmd.AddCommand(collectionList())
	collectionCmd.AddCommand(collectionDelete())
}

func collectionCreate() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func collectionList() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List collections",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func collectionDelete() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a collection",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}
