# Implementation Plan

- [x] 1. Create reusable OpenCV setup action

  - Create `.github/actions/setup-opencv/action.yml` composite action that handles OpenCV installation for both Ubuntu and macOS
  - Implement OS detection and appropriate package installation commands
  - Set up proper CGO environment variables and pkg-config paths
  - Add error handling and fallback mechanisms for OpenCV detection
  - _Requirements: 2.1, 2.2, 2.3_

- [x] 2. Create consolidated CI workflow

  - Create `.github/workflows/ci.yml` that replaces build.yml, lint.yml, and lint-soft.yml
  - Implement matrix strategy for Go versions 1.22 and 1.23 on ubuntu-latest and macos-latest
  - Configure job dependencies and parallel execution for optimal performance
  - _Requirements: 1.1, 1.2, 1.3, 3.1, 3.2, 3.3_

- [x] 3. Implement build and test job in CI workflow

  - Add build-and-test job that uses the OpenCV setup action
  - Configure Go module caching for faster builds
  - Implement proper artifact handling for build outputs
  - Add comprehensive error handling and debugging output
  - _Requirements: 1.1, 2.1, 2.2, 2.3, 7.3_

- [x] 4. Implement strict linting job in CI workflow

  - Add lint-strict job that runs golangci-lint with main configuration
  - Configure the job to fail the workflow on linting issues
  - Set up proper CGO environment for gocv dependencies
  - Implement caching for golangci-lint binary and analysis results
  - _Requirements: 4.1, 4.3, 7.1, 7.3_

- [x] 5. Implement soft linting job in CI workflow

  - Add lint-soft job that runs golangci-lint with soft configuration
  - Configure the job to report issues without failing the workflow
  - Set up conditional execution for pull requests and development branches
  - Use `continue-on-error: true` and appropriate exit codes
  - _Requirements: 4.2, 4.3, 7.4_

- [x] 6. Implement coverage reporting in CI workflow

  - Add coverage generation step to the build-and-test job
  - Configure coverage report generation in both text and HTML formats
  - Set up artifact upload for coverage reports with 30-day retention
  - Add coverage summary output for PR comments
  - _Requirements: 6.1, 6.2, 6.3, 7.3_

- [x] 7. Update GitLab mirror workflow

  - Update `.github/workflows/mirror-to-gitlab.yml` to use actions/checkout@v4
  - Improve error handling and add better logging for mirror operations
  - Add conditional execution based on repository settings and secret availability
  - Implement proper git configuration and force-push safety measures
  - _Requirements: 5.1, 5.3, 7.1, 7.2_

- [x] 8. Update badge management workflow

  - Update `.github/workflows/update-readme-badge.yml` to use actions/checkout@v4
  - Improve badge detection and insertion logic with better regex patterns
  - Add error handling for missing GitLab credentials and repository access
  - Implement conditional execution and proper git commit handling
  - _Requirements: 5.2, 7.1, 7.2_

- [x] 9. Add workflow permissions and security configurations

  - Add minimal required permissions to CI workflow using `permissions` block (currently missing)
  - Verify secret validation and error handling in mirror and badge workflows
  - Ensure all action versions are pinned to specific versions for security
  - Set up appropriate `GITHUB_TOKEN` permissions for each job
  - _Requirements: 7.2, 7.4_

- [x] 10. Optimize workflow performance and caching

  - Add Go build cache to CI workflow for faster builds (currently only has module cache)
  - Implement golangci-lint caching for improved performance
  - Optimize artifact naming to avoid conflicts in matrix builds
  - Add conditional execution for soft linting on pull requests only
  - _Requirements: 7.3, 4.2_

- [x] 11. Test and validate workflow functionality

  - Verify matrix builds work correctly across all Go versions and platforms
  - Test OpenCV setup action reliability on both Ubuntu and macOS runners
  - Validate coverage reporting and artifact uploads work as expected
  - Ensure GitLab mirroring and badge updates function properly with secrets
  - _Requirements: 1.1, 1.2, 1.3, 3.1, 3.2, 3.3, 5.1, 5.2, 6.1, 6.2_

- [x] 11. Remove deprecated workflow files
  - Delete `.github/workflows/build.yml` after verifying CI workflow works
  - Delete `.github/workflows/lint.yml` after confirming lint jobs function
  - Delete `.github/workflows/lint-soft.yml` after validating soft lint integration
  - Update any documentation references to old workflow names
  - _Requirements: 1.1, 4.1, 4.2_
