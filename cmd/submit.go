package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-aws-batch-adapter/internal/batch"
	"github.com/scttfrdmn/ood-aws-batch-adapter/internal/ood"
	"github.com/spf13/cobra"
)

var (
	submitQueue  string
	submitJobDef string
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an OOD job to AWS Batch",
	Long:  "Reads a JSON job spec from stdin and submits it to the specified Batch job queue.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec ood.JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		queue := submitQueue
		if queue == "" {
			queue = spec.Queue
		}
		if queue == "" {
			return fmt.Errorf("--queue is required (or set queue in job spec)")
		}

		jobDef := submitJobDef
		if jobDef == "" {
			jobDef = "ood-job:latest"
		}

		ctx := context.Background()
		client, err := batch.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}

		jobID, err := client.SubmitJob(ctx, spec.JobName, queue, jobDef, spec.Script, spec.Env)
		if err != nil {
			return err
		}

		fmt.Println(jobID)
		return nil
	},
}

func init() {
	submitCmd.Flags().StringVar(&submitQueue, "queue", "", "Batch job queue ARN or name")
	submitCmd.Flags().StringVar(&submitJobDef, "job-definition", "", "Batch job definition name:revision (default: ood-job:latest)")
	rootCmd.AddCommand(submitCmd)
}
