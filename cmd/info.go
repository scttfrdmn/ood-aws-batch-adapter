package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-aws-batch-adapter/internal/batch"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <job-id>",
	Short: "Print full Batch job details as JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := batch.New(ctx, region)
		if err != nil {
			return err
		}
		detail, err := client.DescribeJob(ctx, args[0])
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(detail)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
