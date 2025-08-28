# Requirements Document

## Introduction

This feature consolidates and optimizes the GitHub Actions workflows for the asciicam project. Currently, there are 5 separate workflow files with overlapping functionality, redundant OpenCV setup steps, and inconsistent configurations. The goal is to streamline these into a cohesive, maintainable CI/CD pipeline that follows GitHub Actions best practices while maintaining all existing functionality.

## Requirements

### Requirement 1

**User Story:** As a developer, I want a single comprehensive CI workflow that handles building, testing, and linting, so that I have faster feedback and reduced complexity.

#### Acceptance Criteria

1. WHEN code is pushed or a pull request is created THEN the system SHALL run a unified CI workflow that includes build, test, lint, and coverage steps
2. WHEN the CI workflow runs THEN the system SHALL complete in less time than the current separate workflows combined
3. WHEN the CI workflow fails THEN the system SHALL provide clear indication of which step failed (build, test, lint, or coverage)

### Requirement 2

**User Story:** As a maintainer, I want consistent OpenCV setup across all workflows, so that builds are reliable and maintainable.

#### Acceptance Criteria

1. WHEN any workflow requires OpenCV THEN the system SHALL use a standardized OpenCV installation step
2. WHEN OpenCV is installed THEN the system SHALL properly detect and configure pkg-config paths for both opencv and opencv4
3. WHEN CGO environment variables are set THEN the system SHALL use consistent values across all workflow steps

### Requirement 3

**User Story:** As a developer, I want matrix builds for multiple Go versions and platforms, so that compatibility is ensured across supported environments.

#### Acceptance Criteria

1. WHEN the CI workflow runs THEN the system SHALL test against Go 1.22 and 1.23
2. WHEN the CI workflow runs THEN the system SHALL test on ubuntu-latest and macos-latest
3. WHEN matrix builds run THEN the system SHALL fail the entire workflow if any matrix job fails

### Requirement 4

**User Story:** As a maintainer, I want separate linting configurations (strict and soft) to be preserved, so that code quality standards are maintained while allowing flexibility during development.

#### Acceptance Criteria

1. WHEN a pull request is created THEN the system SHALL run strict linting with the main golangci-lint configuration
2. WHEN code is pushed to non-main branches THEN the system SHALL optionally run soft linting that doesn't fail the build
3. WHEN linting runs THEN the system SHALL only report issues in changed files for pull requests

### Requirement 5

**User Story:** As a project maintainer, I want the GitLab mirroring and badge update workflows to remain functional, so that the project maintains its multi-platform presence.

#### Acceptance Criteria

1. WHEN code is pushed to main THEN the system SHALL mirror the repository to GitLab
2. WHEN the mirror completes successfully THEN the system SHALL update the README with appropriate badges
3. WHEN mirroring fails THEN the system SHALL not affect the main CI workflow

### Requirement 6

**User Story:** As a developer, I want coverage reporting integrated into the main workflow, so that code quality metrics are consistently tracked.

#### Acceptance Criteria

1. WHEN tests run THEN the system SHALL generate coverage reports in both text and HTML formats
2. WHEN coverage is generated THEN the system SHALL upload coverage artifacts for download
3. WHEN coverage drops significantly THEN the system SHALL provide clear visibility into the change

### Requirement 7

**User Story:** As a maintainer, I want workflows to use current GitHub Actions best practices, so that the CI system is secure, efficient, and maintainable.

#### Acceptance Criteria

1. WHEN workflows are defined THEN the system SHALL use the latest stable versions of GitHub Actions
2. WHEN secrets are used THEN the system SHALL follow least-privilege access patterns
3. WHEN caching is implemented THEN the system SHALL cache Go modules and build dependencies appropriately
4. WHEN workflows run THEN the system SHALL use appropriate permissions (read-only by default)
