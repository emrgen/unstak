package cmd

import (
	"github.com/emrgen/unpost"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var spaceCommand = &cobra.Command{
	Use:   "space",
	Short: "space commands",
}

func init() {
	spaceCommand.AddCommand(createSpaceCommand())
	spaceCommand.AddCommand(getSpaceCommand())
	spaceCommand.AddCommand(listSpaceCommand())
	spaceCommand.AddCommand(deleteSpaceCommand())
	spaceCommand.AddCommand(addSpaceMemberCommand())
	spaceCommand.AddCommand(removeSpaceMemberCommand())
	spaceCommand.AddCommand(listSpaceMembersCommand())
	spaceCommand.AddCommand(updateSpaceMemberCommand())
}

func createSpaceCommand() *cobra.Command {
	var spaceName string
	var pool bool

	command := &cobra.Command{
		Use:   "create",
		Short: "create a space",
		Run: func(cmd *cobra.Command, args []string) {
			if spaceName == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			client, err := unpost.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			req := &v1.CreateSpaceRequest{
				Name: spaceName,
			}
			if pool {
				req.PoolName = spaceName
			}

			res, err := client.CreateSpace(tokenContext(), req)
			if err != nil {
				logrus.Errorf("error creating space: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			table.Append([]string{res.Space.Id, res.Space.Name})
			table.Render()
		},
	}

	command.Flags().StringVarP(&spaceName, "name", "n", "", "space name")
	command.Flags().BoolVarP(&pool, "pool", "p", false, "create a pool in authbase project")

	return command
}

func getSpaceCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get a space",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func listSpaceCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list spaces",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := unpost.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			res, err := client.ListSpace(tokenContext(), &v1.ListSpaceRequest{})
			if err != nil {
				logrus.Errorf("error listing spaces: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			for _, space := range res.Spaces {
				table.Append([]string{space.Id, space.Name})
			}
			table.Render()
		},
	}

	return command
}

func deleteSpaceCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete a space",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func addSpaceMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "add-member",
		Short: "add a member to a space",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func removeSpaceMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "remove-member",
		Short: "remove a member from a space",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func listSpaceMembersCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list-members",
		Short: "list members of a space",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func updateSpaceMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update-member",
		Short: "update a member of a space",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}
