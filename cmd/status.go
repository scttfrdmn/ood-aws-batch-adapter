package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-aws-batch-adapter/internal/batch"
	"github.com/scttfrdmn/ood-aws-batch-adapter/internal/ood"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <job-id>",
	Short: "Get the status of a Batch job",
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

		status := batchStateToOod(detail.Status)
		exitCode := 0
		if detail.Container != nil && detail.Container.ExitCode != nil {
			exitCode = int(*detail.Container.ExitCode)
		}

		js := ood.JobStatus{
			ID:       args[0],
			Status:   status,
			ExitCode: exitCode,
		}

		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func batchStateToOod(s types.JobStatus) string {
	switch s {
	case types.JobStatusSubmitted, types.JobStatusPending, types.JobStatusRunnable:
		return ood.StatusQueued
	case types.JobStatusStarting, types.JobStatusRunning:
		return ood.StatusRunning
	case types.JobStatusSucceeded:
		return ood.StatusCompleted
	case types.JobStatusFailed:
		return ood.StatusFailed
	default:
		return ood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// suppress unused import warning for fmt
var _ = fmt.Sprintf
