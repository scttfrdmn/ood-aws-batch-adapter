# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial scaffold — OOD compute adapter for AWS Batch, translating Open OnDemand job submissions to AWS Batch API calls.
- CLI commands: `submit` (job spec as JSON from stdin → AWS Batch), `status <job-id>` (OOD-normalized status), `delete <job-id>` (terminate), and `info <job-id>` (full Batch job details as JSON).
- Unit tests for Batch → OOD status state mapping.
- Substrate integration tests for the Batch job lifecycle.
- CI workflow with pinned action SHAs.
