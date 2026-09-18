# Sign-off Validation Specification

## Purpose

This specification defines the requirements for verifying that Git commits include a Developer Certificate of Origin (DCO) sign-off footer, supporting contributor email whitelisting, branch regex pattern whitelisting, and commit message regex pattern whitelisting.

## Requirements

### Requirement: DCO Sign-off Footer Verification

The system SHALL verify that each commit message in the evaluated commit range contains a valid `Signed-off-by:` footer line.

#### Scenario: Commit contains valid sign-off footer

- **WHEN** a commit message contains a line matching `Signed-off-by: Developer Name <developer@example.com>`
- **THEN** the sign-off validator marks the commit as signed off, increments the assertion count, and reports valid status.

#### Scenario: Commit lacks sign-off footer

- **WHEN** a commit message does not contain a `Signed-off-by:` line
- **THEN** the sign-off validator marks the commit as not signed off, increments the failure count, and includes a failure message identifying the non-compliant commit hash.

### Requirement: Sign-off Contributor Email Whitelisting

The system SHALL skip sign-off validation for commits authored by contributors whose email is configured in `validate.signoff.whitelist`.

#### Scenario: Commit author email is in sign-off whitelist

- **WHEN** a commit's author email matches an entry in `validate.signoff.whitelist`
- **THEN** the sign-off validator skips verification for that commit, increments the skipped count, records the author email as the skip reason, and does not register a failure.

### Requirement: Sign-off Branch Pattern Whitelisting

The system SHALL skip sign-off validation for all commits in the evaluated range when the current branch matches any regular expression pattern in `validate.signoff.branch_patterns`.

#### Scenario: Current branch matches sign-off branch pattern

- **WHEN** `validate.signoff.branch_patterns` contains `^release-please--.*` and the current branch matches
- **THEN** the sign-off validator skips verification for all commits in the range, marks them with a skip reason identifying the matched branch pattern, increments skipped count, and reports valid status.

### Requirement: Sign-off Message Pattern Whitelisting

The system SHALL skip sign-off validation for any commit whose message matches any regular expression pattern in `validate.signoff.patterns`.

#### Scenario: Commit message matches sign-off message pattern

- **WHEN** a commit message is "chore(main): release 0.6.0" without a sign-off footer, and `validate.signoff.patterns` contains `^chore\(.*\): release .*`
- **THEN** the sign-off validator skips verification for that commit, records the matched pattern as the skip reason, increments skipped count, and reports valid status.

#### Scenario: Commit message does not match sign-off message pattern

- **WHEN** an unsigned commit message does not match any pattern in `validate.signoff.patterns` and is not whitelisted by email or branch pattern
- **THEN** the sign-off validator marks the commit as failed due to missing sign-off footer.
