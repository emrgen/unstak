package cmd

import (
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var outletCmd = &cobra.Command{
	Use:   "outlet",
	Short: "outlet commands",
}

func init() {
	outletCmd.AddCommand(outletCreate())
	outletCmd.AddCommand(outletList())
	outletCmd.AddCommand(outletDelete())
	outletCmd.AddCommand(outletAddMember())
	outletCmd.AddCommand(outletUpdateMember())
	outletCmd.AddCommand(outletListMembers())
	outletCmd.AddCommand(outletRemoveMember())
}

func outletCreate() *cobra.Command {
	var outletName string
	var projectID string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a outlet",
		Run: func(cmd *cobra.Command, args []string) {
			client, close := outletClient()
			defer close()

			if outletName == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			res, err := client.CreateOutlet(tokenContext(), &v1.CreateOutletRequest{
				Name:        outletName,
				Description: "",
			},
			)
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Infof("outlet created with id: %s", res.Outlet.Id)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			table.Append([]string{res.Outlet.Id, res.Outlet.Name})
			table.Render()

		},
	}

	command.Flags().StringVarP(&outletName, "name", "n", "", "name of the outlet")
	command.Flags().StringVarP(&projectID, "project", "p", "", "project id to create the outlet in")

	return command
}

func outletList() *cobra.Command {
	var projectID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List outlets",
		Run: func(cmd *cobra.Command, args []string) {
			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
			}

			client, close := outletClient()
			defer close()

			res, err := client.ListOutlets(tokenContext(), &v1.ListOutletsRequest{})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			for _, outlet := range res.Outlets {
				table.Append([]string{outlet.Id, outlet.Name})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectID, "project", "p", "", "project id to list the outlets in")

	return command
}

func outletDelete() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a outlet",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func outletAddMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "add-member",
		Short: "Add a member to a outlet",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func outletUpdateMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "update-member",
		Short: "Update a member of a outlet",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func outletListMembers() *cobra.Command {
	command := &cobra.Command{
		Use:   "list-members",
		Short: "List members of a outlet",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func outletRemoveMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "remove-member",
		Short: "Remove a member from a outlet",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}
