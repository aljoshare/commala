# Commit Message Validation Specification

## Purpose

This specification defines the requirements for validating Git commit messages against the Conventional Commits specification, supporting contributor email whitelisting, commit message regex pattern whitelisting, and branch-level pattern whitelisting.

## Requirements

### Requirement: Conventional Commits Format Enforcement

The system SHALL validate that every commit message in the specified range conforms to the Conventional Commits specification, requiring an allowed type, optional scope, optional breaking change exclamation mark, colon separator, subject description, and optional body/footers.

#### Scenario: Valid conventional commit messages

- **WHEN** a commit message has the form "feat: add user authentication", "fix(api)!: drop deprecated endpoints", or "docs(readme): fix typos\n\nDetailed body explanation"
- **THEN** the message validator marks the commit as conventional, increments assertions, and reports validation success.

#### Scenario: Invalid commit message format

- **WHEN** a commit message does not match conventional commit formatting (such as "update code", "WIP", or "fixed stuff")
- **THEN** the message validator marks the commit as invalid, increments the failure count, and includes a descriptive failure message for the commit hash.

### Requirement: Message Contributor Email Whitelisting

The system SHALL skip commit message validation for any individual commit authored by a contributor whose email is configured in `validate.message.whitelist`.

#### Scenario: Commit author email matches whitelist

- **WHEN** a commit's author email matches an entry in `validate.message.whitelist`
- **THEN** the message validator skips validation for that specific commit, increments the skipped count, records the author email as the skip reason, and does not count it as a failure.

#### Scenario: Non-whitelisted commit author

- **WHEN** a commit's author email is not present in `validate.message.whitelist`
- **THEN** the message validator evaluates the commit message against configured patterns and conventional commit rules.

### Requirement: Message Branch Pattern Whitelisting

The system SHALL skip conventional commit validation for all commits in the evaluated range when the current branch matches any regular expression configured in `validate.message.branch_patterns`.

#### Scenario: Branch matches message branch pattern

- **WHEN** `validate.message.branch_patterns` contains `^release-please--.*` and the current branch name matches
- **THEN** the message validator marks all commits in the range as skipped with a reason referencing the matched branch pattern, increments the skipped count, and reports valid status.

### Requirement: Commit Message Pattern Whitelisting

The system SHALL skip conventional commit validation for any individual commit whose message matches any regular expression configured in `validate.message.patterns`.

#### Scenario: Commit message matches message pattern

- **WHEN** a commit message is "chore(main): release 0.6.0" and `validate.message.patterns` contains `^chore\(.*\): release .*`
- **THEN** the message validator skips validation for that commit, records the matching message pattern as the skip reason, increments skipped count, and treats the commit as valid.

#### Scenario: Commit message does not match message pattern

- **WHEN** a commit message does not match any pattern in `validate.message.patterns` and is not whitelisted by email or branch pattern
- **THEN** the message validator executes standard conventional commit format validation on the commit message.
