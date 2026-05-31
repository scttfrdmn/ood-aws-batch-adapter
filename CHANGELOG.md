# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-30

### Fixed
- Repinned `github.com/scttfrdmn/substrate` 0.45.2 → 0.65.0 and regenerated go.sum. The v0.45.2 tag content was changed upstream (substrate#296), so the recorded checksum no longer matched and `go test -tags=integration` failed with a go.sum SECURITY ERROR. Integration tests now build and pass.

### Added
- Initial scaffold — OOD compute adapter for AWS Batch, translating Open OnDemand job submissions to AWS Batch API calls.
- CLI commands: `submit` (job spec as JSON from stdin → AWS Batch), `status <job-id>` (OOD-normalized status), `delete <job-id>` (terminate), and `info <job-id>` (full Batch job details as JSON).
- Unit tests for Batch → OOD status state mapping.
- Substrate integration tests for the Batch job lifecycle.
- CI workflow with pinned action SHAs.
