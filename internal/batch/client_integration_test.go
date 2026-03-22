//go:build integration

package batch_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsbatch "github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	substrate "github.com/scttfrdmn/substrate"

	. "github.com/scttfrdmn/ood-aws-batch-adapter/internal/batch"
)

// substrateBatchClient returns a raw AWS Batch SDK client pointed at the substrate server.
func substrateBatchClient(t *testing.T, endpointURL string) *awsbatch.Client {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(endpointURL),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	return awsbatch.NewFromConfig(cfg)
}

// setupBatchPrereqs creates the prerequisite Batch resources (compute environment,
// job queue, job definition) and returns the queue and job definition ARNs.
func setupBatchPrereqs(t *testing.T, ctx context.Context, raw *awsbatch.Client) (queueARN, jobDefARN string) {
	t.Helper()

	ceOut, err := raw.CreateComputeEnvironment(ctx, &awsbatch.CreateComputeEnvironmentInput{
		ComputeEnvironmentName: aws.String("ood-test-ce"),
		Type:                   "MANAGED",
	})
	if err != nil {
		t.Fatalf("CreateComputeEnvironment: %v", err)
	}
	ceARN := aws.ToString(ceOut.ComputeEnvironmentArn)
	t.Logf("compute environment: %s", ceARN)

	jqOut, err := raw.CreateJobQueue(ctx, &awsbatch.CreateJobQueueInput{
		JobQueueName: aws.String("ood-test-queue"),
		Priority:     aws.Int32(1),
		ComputeEnvironmentOrder: []batchtypes.ComputeEnvironmentOrder{
			{ComputeEnvironment: aws.String(ceARN), Order: aws.Int32(1)},
		},
	})
	if err != nil {
		t.Fatalf("CreateJobQueue: %v", err)
	}
	queueARN = aws.ToString(jqOut.JobQueueArn)
	t.Logf("job queue: %s", queueARN)

	jdOut, err := raw.RegisterJobDefinition(ctx, &awsbatch.RegisterJobDefinitionInput{
		JobDefinitionName: aws.String("ood-test-jobdef"),
		Type:              "container",
	})
	if err != nil {
		t.Fatalf("RegisterJobDefinition: %v", err)
	}
	jobDefARN = aws.ToString(jdOut.JobDefinitionArn)
	t.Logf("job definition: %s", jobDefARN)

	return queueARN, jobDefARN
}

// TestSubmitDescribeTerminateJob_Substrate exercises the full Batch job lifecycle
// via the adapter client against the substrate emulator.
func TestSubmitDescribeTerminateJob_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	raw := substrateBatchClient(t, ts.URL)
	queueARN, jobDefARN := setupBatchPrereqs(t, ctx, raw)

	// Build the adapter client — it picks up AWS_ENDPOINT_URL from the environment.
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// SubmitJob
	jobID, err := client.SubmitJob(ctx,
		"ood-integration-test",
		queueARN,
		jobDefARN,
		"#!/bin/bash\necho hello",
		map[string]string{"TEST_VAR": "hello"},
	)
	if err != nil {
		t.Fatalf("SubmitJob: %v", err)
	}
	if jobID == "" {
		t.Fatal("expected non-empty job ID")
	}
	t.Logf("submitted job ID: %s", jobID)

	// DescribeJob
	detail, err := client.DescribeJob(ctx, jobID)
	if err != nil {
		t.Fatalf("DescribeJob: %v", err)
	}
	if aws.ToString(detail.JobId) != jobID {
		t.Errorf("DescribeJob: got ID %q, want %q", aws.ToString(detail.JobId), jobID)
	}
	t.Logf("job status: %s", detail.Status)

	// TerminateJob
	err = client.TerminateJob(ctx, jobID, "integration-test-cleanup")
	if err != nil {
		t.Fatalf("TerminateJob: %v", err)
	}
	t.Log("job terminated successfully")
}

// TestDescribeJob_NotFound_Substrate verifies that DescribeJob returns an error
// when the job ID does not exist.
func TestDescribeJob_NotFound_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.DescribeJob(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error for non-existent job, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Logf("error (acceptable): %v", err)
	}
}
