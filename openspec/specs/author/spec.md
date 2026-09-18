# Author Validation Specification

## Purpose

This specification defines the requirements for verifying that Git commits have non-empty author names and non-empty author email addresses, supporting contributor email whitelisting for both checks.

## Requirements

### Requirement: Author Name Presence Verification

The system SHALL verify that every commit in the evaluated commit range has a non-empty author name.

#### Scenario: Commit has valid author name

- **WHEN** a commit contains a non-empty author name
- **THEN** the author name validator marks the commit as valid, increments assertions, and reports success.

#### Scenario: Commit has missing author name

- **WHEN** a commit has an empty author name string
- **THEN** the author name validator marks the commit as invalid, increments the failure count, and includes a descriptive failure message for the commit hash.

### Requirement: Author Email Presence Verification

The system SHALL verify that every commit in the evaluated commit range has a non-empty author email address.

#### Scenario: Commit has valid author email

- **WHEN** a commit contains a non-empty author email address
- **THEN** the author email validator marks the commit as valid, increments assertions, and reports success.

#### Scenario: Commit has missing author email

- **WHEN** a commit has an empty author email string
- **THEN** the author email validator marks the commit as invalid, increments the failure count, and includes a descriptive failure message for the commit hash.

### Requirement: Author Validation Contributor Whitelisting

The system SHALL skip author name or email validation for commits whose author email is listed in `validate.author.name.whitelist` or `validate.author.email.whitelist`.

#### Scenario: Commit author email matches author name whitelist

- **WHEN** a commit's author email matches an entry in `validate.author.name.whitelist`
- **THEN** the author name validator skips validation for that commit, increments the skipped count, records the whitelisted email, and reports valid status.

#### Scenario: Commit author email matches author email whitelist

- **WHEN** a commit's author email matches an entry in `validate.author.email.whitelist`
- **THEN** the author email validator skips validation for that commit, increments the skipped count, records the whitelisted email, and reports valid status.
