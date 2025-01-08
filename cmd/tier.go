package cmd

import (
	"github.com/emrgen/unpost"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var tierCmd = &cobra.Command{
	Use:   "tier",
	Short: "tier commands",
}

func init() {
	tierCmd.AddCommand(tierCreate())
	tierCmd.AddCommand(tierList())
	tierCmd.AddCommand(tierDelete())
	tierCmd.AddCommand(tierAddMember())
	tierCmd.AddCommand(tierUpdateMember())
	tierCmd.AddCommand(tierListMembers())
	tierCmd.AddCommand(tierRemoveMember())
}

func tierCreate() *cobra.Command {
	var tierName string
	var projectID string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a tier",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			if tierName == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			res, err := client.CreateTier(tokenContext(), &v1.CreateTierRequest{
				Name:        tierName,
				Description: "",
			},
			)
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Infof("tier created with id: %s", res.Tier.Id)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			table.Append([]string{res.Tier.Id, res.Tier.Name})
			table.Render()

		},
	}

	command.Flags().StringVarP(&tierName, "name", "n", "", "name of the tier")
	command.Flags().StringVarP(&projectID, "project", "p", "", "project id to create the tier in")

	return command
}

func tierList() *cobra.Command {
	var projectID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List tiers",
		Run: func(cmd *cobra.Command, args []string) {
			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.ListTiers(tokenContext(), &v1.ListTiersRequest{})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			for _, tier := range res.Tiers {
				table.Append([]string{tier.Id, tier.Name})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectID, "project", "p", "", "project id to list the tiers in")

	return command
}

func tierDelete() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tier",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func tierAddMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "add-member",
		Short: "Add a member to a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func tierUpdateMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "update-member",
		Short: "Update a member of a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func tierListMembers() *cobra.Command {
	command := &cobra.Command{
		Use:   "list-members",
		Short: "List members of a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func tierRemoveMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "remove-member",
		Short: "Remove a member from a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}
