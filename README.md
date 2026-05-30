# ood-aws-batch-adapter

[![CI](https://github.com/scttfrdmn/ood-aws-batch-adapter/actions/workflows/ci.yml/badge.svg)](https://github.com/scttfrdmn/ood-aws-batch-adapter/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/ood-aws-batch-adapter)](https://goreportcard.com/report/github.com/scttfrdmn/ood-aws-batch-adapter)
[![codecov](https://codecov.io/gh/scttfrdmn/ood-aws-batch-adapter/branch/main/graph/badge.svg)](https://codecov.io/gh/scttfrdmn/ood-aws-batch-adapter)
[![Go Reference](https://pkg.go.dev/badge/github.com/scttfrdmn/ood-aws-batch-adapter.svg)](https://pkg.go.dev/github.com/scttfrdmn/ood-aws-batch-adapter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

OOD compute adapter for AWS Batch. Translates Open OnDemand job submissions to AWS Batch API calls.

## Commands

| Command | Description |
|---------|-------------|
| `submit` | Submit an OOD job spec (JSON from stdin) to AWS Batch |
| `status <job-id>` | Get OOD-normalized status of a Batch job |
| `delete <job-id>` | Terminate a Batch job |
| `info <job-id>` | Print full Batch job details as JSON |

## Usage

```bash
# Submit a job (OOD job spec on stdin)
echo '{"job_name":"myjob","script":"#!/bin/bash\necho hello"}' | \
  ood-aws-batch-adapter submit --queue ood-test --region us-east-1

# Check status
ood-aws-batch-adapter status <job-id>

# Terminate
ood-aws-batch-adapter delete <job-id>
```

## OOD Cluster Config

```yaml
# /etc/ood/config/clusters.d/aws-batch.yml
---
v2:
  metadata:
    title: "AWS Batch"
  job:
    adapter: "adapter_script"
    submit_host: "localhost"
    submit:
      script: "/usr/local/lib/ood-adapters/ood-aws-batch-adapter"
      args:
        - submit
        - "--queue=arn:aws:batch:us-east-1:123456789012:job-queue/ood-test"
        - "--region=us-east-1"
```

## Infrastructure

Terraform in `aws-openondemand` with `adapters_enabled = ["batch"]` provisions:
- AWS Batch compute environment (SPOT, up to 256 vCPUs)
- Job queue `ood-<env>`
- IAM policy on the OOD instance role
