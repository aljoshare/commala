# Whitelisting Specification

## Purpose

This specification defines the requirements for global branch-level pattern whitelisting and regular expression compilation to bypass all validation checks for automated workflows, merge requests, and bot accounts.

## Requirements

### Requirement: Global Branch Regex Whitelisting

The system SHALL evaluate the current branch name against the list of regular expression patterns configured in `whitelist.branches` prior to executing any individual validator.

#### Scenario: Branch matches global branch pattern

- **WHEN** the current branch name matches any pattern defined in `whitelist.branches` (e.g., `^release-please--.*` or `^dependabot/.*`)
- **THEN** the system bypasses execution of all enabled validators and returns skipped validation results for all checks.

#### Scenario: Branch does not match global branch pattern

- **WHEN** the current branch does not match any pattern in `whitelist.branches`
- **THEN** the system proceeds to execute all enabled validators concurrently.

### Requirement: Bypass All Validations on Global Match

The system SHALL generate skipped result records for every enabled validator when a global branch pattern matches, ensuring all checks are recorded as skipped without triggering failures or non-zero exit codes.

#### Scenario: Result generation on global branch whitelist match

- **WHEN** a global branch pattern matches for a commit range
- **THEN** each enabled validator (Author Name, Author Email, Branch, Conventional Message, Signed-Off Message) generates a result with valid status true, failure count 0, skipped count matching the commits (or 1 for branch), and a skip reason stating "Branch whitelisted: <pattern>".

#### Scenario: Exit code on global branch whitelist match

- **WHEN** all validations are skipped due to a global branch whitelist match
- **THEN** the application exits with status code 0.

### Requirement: Regular Expression Compilation and Error Handling

The system SHALL compile all configured branch and message patterns using standard Go regular expression syntax, failing fast with a descriptive error if any pattern is malformed.

#### Scenario: Valid regular expressions

- **WHEN** patterns contain valid Go regular expressions such as `^release-please--.*` or `^chore\(.*\): release .*`
- **THEN** the patterns compile successfully and match target strings according to regex rules.

#### Scenario: Malformed regular expression syntax

- **WHEN** a configured pattern contains invalid regex syntax (such as an unclosed parenthesis `[a-z`)
- **THEN** pattern compilation returns an error and execution terminates with a non-zero exit code.
