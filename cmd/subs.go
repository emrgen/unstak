package cmd

import (
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "subscription commands",
}

func init() {
	subscriptionCmd.AddCommand(subscriptionCreate())
	subscriptionCmd.AddCommand(subscriptionList())
	subscriptionCmd.AddCommand(subscriptionDelete())
	subscriptionCmd.AddCommand(subscriptionAddMember())
	subscriptionCmd.AddCommand(subscriptionUpdateMember())
	subscriptionCmd.AddCommand(subscriptionListMembers())
	subscriptionCmd.AddCommand(subscriptionRemoveMember())
}

func subscriptionCreate() *cobra.Command {
	var subscriptionName string
	var projectID string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a subscription",
		Run: func(cmd *cobra.Command, args []string) {
			client, close := subscriptionClient()
			defer close()

			if subscriptionName == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			res, err := client.CreateSubscription(tokenContext(), &v1.CreateSubscriptionRequest{
				Name:        subscriptionName,
				Description: "",
			},
			)
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Infof("subscription created with id: %s", res.Subscription.Id)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			table.Append([]string{res.Subscription.Id, res.Subscription.Name})
			table.Render()

		},
	}

	command.Flags().StringVarP(&subscriptionName, "name", "n", "", "name of the subscription")
	command.Flags().StringVarP(&projectID, "project", "p", "", "project id to create the subscription in")

	return command
}

func subscriptionList() *cobra.Command {
	var projectID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List subscriptions",
		Run: func(cmd *cobra.Command, args []string) {
			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
			}

			client, close := subscriptionClient()
			defer close()

			res, err := client.ListSubscriptions(tokenContext(), &v1.ListSubscriptionsRequest{})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name"})
			for _, subscription := range res.Subscriptions {
				table.Append([]string{subscription.Id, subscription.Name})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectID, "project", "p", "", "project id to list the subscriptions in")

	return command
}

func subscriptionDelete() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func subscriptionAddMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "add-member",
		Short: "Add a member to a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func subscriptionUpdateMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "update-member",
		Short: "Update a member of a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func subscriptionListMembers() *cobra.Command {
	command := &cobra.Command{
		Use:   "list-members",
		Short: "List members of a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}

func subscriptionRemoveMember() *cobra.Command {
	command := &cobra.Command{
		Use:   "remove-member",
		Short: "Remove a member from a subscription",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return command
}
