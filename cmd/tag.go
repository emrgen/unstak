package cmd

import (
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag commands",
}

func init() {
	tagCmd.AddCommand(tagCreate())
	tagCmd.AddCommand(tagList())
}

func tagCreate() *cobra.Command {
	var tagName string
	var spaceID string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a tag",
		Run: func(cmd *cobra.Command, args []string) {
			if spaceID == "" {
				cmd.Println("space id is required")
				return
			}

			if tagName == "" {
				cmd.Println("tag name is required")
				return
			}

			client, close := tagClient()
			defer close()

			res, err := client.CreateTag(tokenContext(), &v1.CreateTagRequest{
				SpaceId: spaceID,
				Name:    tagName,
			})
			if err != nil {
				cmd.Println(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			table.Append([]string{res.Tag.Id, res.Tag.Name})
			table.Render()
		},
	}

	command.Flags().StringVarP(&spaceID, "space", "s", "", "space id")
	command.Flags().StringVarP(&tagName, "name", "n", "", "tag name")

	return command
}

func tagList() *cobra.Command {
	var spaceID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Run: func(cmd *cobra.Command, args []string) {
			if spaceID == "" {
				cmd.Println("space id is required")
				return
			}

			client, close := tagClient()
			defer close()

			res, err := client.ListTag(tokenContext(), &v1.ListTagRequest{
				SpaceId: spaceID,
			})
			if err != nil {
				cmd.Println(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})

			for _, tag := range res.Tags {
				table.Append([]string{tag.Id, tag.Name})
			}

			table.Render()
		},
	}

	command.Flags().StringVarP(&spaceID, "space", "s", "", "space id")

	return command
}
