// Package batch wraps the AWS Batch API for the OOD adapter.
package batch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
)

// Client wraps the AWS Batch client.
type Client struct {
	svc    *batch.Client
	region string
}

// New creates a Batch client using the default AWS credential chain.
func New(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: batch.NewFromConfig(cfg), region: region}, nil
}

// SubmitJob submits a containerized job to the given queue.
func (c *Client) SubmitJob(ctx context.Context, jobName, jobQueue, jobDef, script string, env map[string]string) (string, error) {
	var envVars []types.KeyValuePair
	for k, v := range env {
		envVars = append(envVars, types.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}
	envVars = append(envVars, types.KeyValuePair{
		Name:  aws.String("OOD_JOB_SCRIPT"),
		Value: aws.String(script),
	})

	out, err := c.svc.SubmitJob(ctx, &batch.SubmitJobInput{
		JobName:       aws.String(jobName),
		JobQueue:      aws.String(jobQueue),
		JobDefinition: aws.String(jobDef),
		ContainerOverrides: &types.ContainerOverrides{
			Environment: envVars,
		},
	})
	if err != nil {
		return "", fmt.Errorf("batch SubmitJob: %w", err)
	}
	return aws.ToString(out.JobId), nil
}

// DescribeJob returns the current status of a Batch job.
func (c *Client) DescribeJob(ctx context.Context, jobID string) (*types.JobDetail, error) {
	out, err := c.svc.DescribeJobs(ctx, &batch.DescribeJobsInput{
		Jobs: []string{jobID},
	})
	if err != nil {
		return nil, fmt.Errorf("batch DescribeJobs: %w", err)
	}
	if len(out.Jobs) == 0 {
		return nil, fmt.Errorf("job %q not found", jobID)
	}
	return &out.Jobs[0], nil
}

// TerminateJob cancels or terminates a Batch job.
func (c *Client) TerminateJob(ctx context.Context, jobID, reason string) error {
	_, err := c.svc.TerminateJob(ctx, &batch.TerminateJobInput{
		JobId:  aws.String(jobID),
		Reason: aws.String(reason),
	})
	if err != nil {
		return fmt.Errorf("batch TerminateJob: %w", err)
	}
	return nil
}
