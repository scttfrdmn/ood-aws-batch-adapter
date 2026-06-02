package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-aws-batch-adapter/internal/batch"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <job-id>",
	Short: "Terminate a Batch job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := batch.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}
		if err := client.TerminateJob(ctx, args[0], "Cancelled via OOD"); err != nil {
			return err
		}
		fmt.Printf("Job %s terminated\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
