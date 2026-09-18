# Branch Validation Specification

## Purpose

This specification defines the requirements for validating Git branch names against the Conventional Branch specification, supporting contributor email whitelisting, pattern-based regex whitelisting, CI environment detached HEAD detection, and manual branch name overrides.

## Requirements

### Requirement: Conventional Branch Naming Enforcement

The system SHALL validate that git branch names conform to the Conventional Branch naming convention, requiring allowed prefixes followed by a forward slash and a valid branch suffix.

#### Scenario: Valid conventional branch prefix and name

- **WHEN** the current branch name is "feature/login-page", "feat/auth", "bugfix/issue-12", "fix/crash", "hotfix/urgent-patch", "release/v1.0.0", "chore/update-deps", "docs/readme", "test/unit-tests", or "refactor/cleanup"
- **THEN** the branch validator reports valid with an assertion count of 1 and 0 failures.

#### Scenario: Invalid branch naming convention

- **WHEN** the current branch name does not match the allowed pattern (such as "random-branch", "my_feature", or "issue-123")
- **THEN** the branch validator reports invalid, increments the failure count, and includes a descriptive message indicating the branch name is not conventional.

### Requirement: Branch Contributor Email Whitelisting

The system SHALL skip branch name validation when the branch author's email matches an address configured in the contributor whitelist.

#### Scenario: Branch author email is in whitelist

- **WHEN** the author email of the latest commit on the branch matches an entry in `validate.branch.whitelist`
- **THEN** the branch validator skips validation, increments the skipped count, records the whitelisted email in the result message, and marks the result as valid.

#### Scenario: Branch author email is not in whitelist

- **WHEN** the branch author email does not match any entry in `validate.branch.whitelist`
- **THEN** the branch validator proceeds to evaluate the branch name against configured patterns and conventions.

### Requirement: Branch Pattern Whitelisting

The system SHALL skip branch name validation when the branch name matches any regular expression pattern configured in `validate.branch.patterns`.

#### Scenario: Branch name matches whitelisted regex pattern

- **WHEN** the branch name is "release-please--branches--main" and `validate.branch.patterns` contains `^release-please--.*`
- **THEN** the branch validator skips validation, increments the skipped count, marks the result as valid, and records a skip reason indicating the matching branch pattern.

#### Scenario: Branch name does not match whitelisted regex pattern

- **WHEN** the branch name does not match any pattern in `validate.branch.patterns` and is not whitelisted by email
- **THEN** the branch validator performs conventional branch format validation.

### Requirement: CI Detached HEAD Branch Detection

The system SHALL automatically detect the target branch name in detached HEAD CI environments by inspecting standard CI provider environment variables.

#### Scenario: Detection via GitHub Actions environment

- **WHEN** the Git repository is in a detached HEAD state and `GITHUB_HEAD_REF` is set to "feature/ci-build"
- **THEN** the system resolves the branch name as "feature/ci-build" for validation.

#### Scenario: Detection via GitLab CI environment

- **WHEN** the repository is in a detached HEAD state and `CI_MERGE_REQUEST_SOURCE_BRANCH_NAME` (or `CI_COMMIT_REF_NAME` or `CI_COMMIT_BRANCH`) is set
- **THEN** the system resolves the branch name from the respective GitLab CI environment variable.

#### Scenario: Detection via Bitbucket Pipelines environment

- **WHEN** the repository is in a detached HEAD state and `BITBUCKET_BRANCH` is set
- **THEN** the system resolves the branch name from `BITBUCKET_BRANCH`.

### Requirement: Manual Branch Name Override

The system SHALL allow users to explicitly specify or override the detected branch name via CLI flags or environment variables.

#### Scenario: Branch name override via flag or environment variable

- **WHEN** the `--branch-name` flag or `COMMALA_BRANCH_NAME` environment variable is provided with value "feature/manual-override"
- **THEN** the system uses "feature/manual-override" as the branch name for all branch evaluations regardless of Git HEAD or CI environment variables.
