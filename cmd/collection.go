package cmd

import (
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
	var name string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			client, close := collectionClient()
			defer close()

			res, err := client.CreateCollection(tokenContext(), &v1.CreateCollectionRequest{Name: name})
			if err != nil {
				logrus.Errorf("Error creating collection: %v", err)
				return
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Name", "Created At", "Updated At"})
			table.Append([]string{
				res.Collection.Id,
				res.Collection.Name,
				res.Collection.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
				res.Collection.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "Collection name")

	return command
}

func collectionList() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List collections",
		Run: func(cmd *cobra.Command, args []string) {
			client, close := collectionClient()
			defer close()

			res, err := client.ListCollection(tokenContext(), &v1.ListCollectionRequest{})
			if err != nil {
				logrus.Errorf("Error creating collection: %v", err)
				return
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Name", "Created At", "Updated At"})
			for _, collection := range res.Collections {
				table.Append([]string{
					collection.Id,
					collection.Name,
					collection.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
					collection.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"),
				})
			}
			table.Render()
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
