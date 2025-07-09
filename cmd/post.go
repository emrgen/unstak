package cmd

import (
	"fmt"
	"github.com/emrgen/unpost"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "post commands",
}

func init() {
	postCmd.AddCommand(postCreate())
	postCmd.AddCommand(getPost())
	postCmd.AddCommand(postList())
	postCmd.AddCommand(updatePost())
	postCmd.AddCommand(deletePost())
	postCmd.AddCommand(addPostTag())
	postCmd.AddCommand(removePostTag())
	postCmd.AddCommand(updatePostStatus())
}

func postCreate() *cobra.Command {
	var title string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a post",
		Run: func(cmd *cobra.Command, args []string) {
			if title == "" {
				logrus.Errorf("missing required flag: --title")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			res, err := client.CreatePost(tokenContext(), &v1.CreatePostRequest{
				Title: title,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			cmd.Println("Post created")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Title", "Created At", "Updated At"})
			table.Append([]string{res.Post.Id, res.Post.Title, res.Post.CreatedAt.AsTime().Format("2006-01-02 15:04:05"), res.Post.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&title, "title", "t", "", "title of the post")

	return command
}

func getPost() *cobra.Command {
	var postID string

	command := &cobra.Command{
		Use:   "get",
		Short: "Get a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == "" {
				logrus.Errorf("missing required flag: --post-id")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			res, err := client.GetPost(tokenContext(), &v1.GetPostRequest{
				Id: postID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			fmt.Printf("ID: %s\n", res.Post.Id)
			fmt.Printf("Version: %d\n", res.Post.GetVersion())
			fmt.Printf("Status: %s\n", res.Post.GetStatus().String())
			tags := make([]string, 0)
			for _, tag := range res.Post.Tags {
				tags = append(tags, tag.Name)
			}
			fmt.Printf("Tags: %s\n", strings.Join(tags, ", "))
			fmt.Printf("Title: %s\n", res.Post.GetTitle())
			fmt.Printf("Excerpt: %s\n", res.Post.GetExcerpt())
			fmt.Printf("Summary: %s\n", res.Post.GetSummary())
			fmt.Printf("Content: %s\n", res.Post.GetContent())
		},
	}

	command.Flags().StringVarP(&postID, "post-id", "p", "", "post id")

	return command

}

func postList() *cobra.Command {
	var postStatus string
	var spaceID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List posts",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			req := &v1.ListPostRequest{}

			if postStatus != "" {
				var status v1.PostStatus
				switch postStatus {
				case "draft":
					status = v1.PostStatus_DRAFT
				case "published":
					status = v1.PostStatus_PUBLISHED
				case "archived":
					status = v1.PostStatus_ARCHIVED
				}
				req.Status = &status
			}

			res, err := client.ListPost(tokenContext(), req)
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Title", "Created At", "Updated At", "Status", "Version"})
			for _, post := range res.Posts {
				table.Append([]string{post.Id, post.Title, post.CreatedAt.AsTime().Format("2006-01-02 15:04:05"), post.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"), post.Status.String(), fmt.Sprintf("%d", post.GetVersion())})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&spaceID, "space-id", "s", "", "space id")
	command.Flags().StringVarP(&postStatus, "status", "t", "", "status of the post")

	return command
}

func updatePost() *cobra.Command {
	var postID, postTitle, postContent, postSummary, postExcerpt, thumbnail string
	var version int64

	command := &cobra.Command{
		Use:   "update",
		Short: "Update a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == "" {
				logrus.Errorf("missing required flag: --post-id")
				return
			}
			if version == 0 {
				logrus.Errorf("missing required flag: --version")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			req := &v1.UpdatePostRequest{
				PostId:  postID,
				Version: version,
			}

			if postTitle != "" {
				req.Title = &postTitle
			}

			if postContent != "" {
				req.Content = &postContent
			}

			if postSummary != "" {
				req.Summary = &postSummary
			}

			if postExcerpt != "" {
				req.Excerpt = &postExcerpt
			}

			if thumbnail != "" {
				req.Thumbnail = &thumbnail
			}

			_, err = client.UpdatePost(tokenContext(), req)
			if err != nil {
				logrus.Error(err)
				return
			}

			cmd.Println("Post updated")
		},
	}

	command.Flags().StringVarP(&postID, "post-id", "p", "", "post id")
	command.Flags().StringVarP(&postTitle, "title", "t", "", "title of the post")
	command.Flags().StringVarP(&postContent, "content", "c", "", "content of the post")
	command.Flags().StringVarP(&postSummary, "summary", "s", "", "summary of the post")
	command.Flags().StringVarP(&postExcerpt, "excerpt", "e", "", "excerpt of the post")
	command.Flags().Int64VarP(&version, "version", "v", 0, "version of the post")
	command.Flags().StringVarP(&thumbnail, "thumbnail", "i", "", "thumbnail of the post")

	return command
}

func deletePost() *cobra.Command {
	var postID string

	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == "" {
				logrus.Errorf("missing required flag: --post-id")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			_, err = client.DeletePost(tokenContext(), &v1.DeletePostRequest{
				Id: postID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			cmd.Println("Post deleted")

		},
	}

	command.Flags().StringVarP(&postID, "post-id", "p", "", "post id")

	return command
}

func addPostTag() *cobra.Command {
	var postID string
	var tagID string

	command := &cobra.Command{
		Use:   "add-tag",
		Short: "Add a tag to a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == "" {
				logrus.Errorf("missing required flag: --post-id")
				return
			}

			if tagID == "" {
				logrus.Errorf("missing required flag: --tag-id")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			_, err = client.AddPostTag(tokenContext(), &v1.AddPostTagRequest{
				PostId: postID,
				TagId:  tagID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

		},
	}

	command.Flags().StringVarP(&postID, "post-id", "p", "", "post id")
	command.Flags().StringVarP(&tagID, "tag-id", "t", "", "tag id")

	return command
}

func removePostTag() *cobra.Command {
	var postID string
	var tagID string

	command := &cobra.Command{
		Use:   "remove-tag",
		Short: "Remove a tag from a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == "" {
				logrus.Errorf("missing required flag: --post-id")
				return
			}

			if tagID == "" {
				logrus.Errorf("missing required flag: --tag-id")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			_, err = client.RemovePostTag(tokenContext(), &v1.RemovePostTagRequest{
				PostId: postID,
				TagId:  tagID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

		},
	}

	command.Flags().StringVarP(&postID, "post-id", "p", "", "post id")
	command.Flags().StringVarP(&tagID, "tag-id", "t", "", "tag id")

	return command
}

func updatePostStatus() *cobra.Command {
	var postID string
	var status string
	command := &cobra.Command{
		Use:   "status",
		Short: "Publish a post",
		Run: func(cmd *cobra.Command, args []string) {
			if postID == "" {
				logrus.Errorf("missing required flag: --post-id")
				return
			}

			if status == "" {
				logrus.Errorf("missing required flag: --status")
				return
			}

			client, err := unpost.NewClient("8030")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			var postStatus v1.PostStatus
			switch status {
			case "draft":
				postStatus = v1.PostStatus_DRAFT
			case "published":
				postStatus = v1.PostStatus_PUBLISHED
			case "archived":
				postStatus = v1.PostStatus_ARCHIVED
			default:
				logrus.Errorf("invalid status, must be one of draft, published, archived")
			}

			_, err = client.UpdatePostStatus(tokenContext(), &v1.UpdatePostStatusRequest{
				PostId: postID,
				Status: postStatus,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			cmd.Println("Post status updated to", status)
		},
	}

	command.Flags().StringVarP(&postID, "post-id", "p", "", "post id")
	command.Flags().StringVarP(&status, "status", "s", "", "status of the post, one of draft, published, archived")

	return command
}
